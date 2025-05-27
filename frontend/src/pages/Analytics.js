import React, { useState, useEffect } from 'react';
import {
  Container,
  Grid,
  Paper,
  Typography,
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from '@mui/material';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import axios from 'axios';

function Analytics() {
  const [averageGradesByClass, setAverageGradesByClass] = useState([]);
  const [failingStudents, setFailingStudents] = useState([]);
  const [topAndWorstClasses, setTopAndWorstClasses] = useState({
    top_class: '',
    worst_class: '',
  });

  const fetchAnalyticsData = async () => {
    try {
      console.log('Загрузка данных аналитики...');
      const token = localStorage.getItem('token');
      const headers = { Authorization: `Bearer ${token}` };

      const [averageGradesRes, failingStudentsRes, topClassesRes] = await Promise.all([
        axios.get('http://localhost:8000/stats/average-grades', { headers }),
        axios.get('http://localhost:8000/stats/failing-students', { headers }),
        axios.get('http://localhost:8000/stats/top-worst-classes', { headers }),
      ]);

      console.log('Получены данные отстающих студентов:', failingStudentsRes.data);
      console.log('Количество отстающих студентов:', failingStudentsRes.data.length);

      // Преобразуем данные для графика
      const chartData = Object.entries(averageGradesRes.data).map(([className, subjects]) => ({
        class_name: className,
        ...subjects
      }));
      setAverageGradesByClass(chartData);
      setFailingStudents(failingStudentsRes.data);
      setTopAndWorstClasses(topClassesRes.data);
    } catch (error) {
      console.error('Ошибка при получении данных аналитики:', error);
      console.error('Детали ошибки:', error.response?.data);
    }
  };

  // Загружаем данные при монтировании компонента
  useEffect(() => {
    fetchAnalyticsData();
  }, []);

  // Добавляем интервал обновления данных
  useEffect(() => {
    const interval = setInterval(() => {
      fetchAnalyticsData();
    }, 35000); // Обновляем каждые 5 секунд

    return () => clearInterval(interval);
  }, []);

  // Получаем список всех предметов для заголовков таблицы
  const getAllSubjects = () => {
    const subjects = new Set();
    averageGradesByClass.forEach(classData => {
      Object.keys(classData).forEach(key => {
        if (key !== 'class_name') {
          subjects.add(key);
        }
      });
    });
    return Array.from(subjects).sort();
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Typography component="h1" variant="h4" color="primary" gutterBottom>
        Аналитика
      </Typography>

      <Grid container spacing={3}>
        {/* График средних оценок по классам */}
        <Grid item xs={12}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
            <Typography variant="h6" gutterBottom>
              Средние оценки по классам
            </Typography>
            <Box sx={{ height: 400 }}>
              <ResponsiveContainer>
                <BarChart data={averageGradesByClass}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="class_name" />
                  <YAxis domain={[0, 5]} />
                  <Tooltip />
                  <Legend />
                  {getAllSubjects().map((subject, index) => (
                    <Bar
                      key={subject}
                      dataKey={subject}
                      name={subject}
                      fill={`hsl(${index * 45}, 70%, 50%)`}
                    />
                  ))}
                </BarChart>
              </ResponsiveContainer>
            </Box>
          </Paper>
        </Grid>

        {/* Лучший и худший классы */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column', height: 400 }}>
            <Typography variant="h6" gutterBottom>
              Лучший и худший классы
            </Typography>
            <Box sx={{ mt: 2 }}>
              <Typography variant="subtitle1" gutterBottom>
                Лучший класс: {topAndWorstClasses.top_class}
              </Typography>
              <Typography variant="subtitle1">
                Худший класс: {topAndWorstClasses.worst_class}
              </Typography>
            </Box>
          </Paper>
        </Grid>

        {/* Список отстающих студентов */}
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
            <Typography variant="h6" gutterBottom>
              Отстающие студенты
            </Typography>
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>ФИО</TableCell>
                    <TableCell>Класс</TableCell>
                    <TableCell>Предметы</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {failingStudents.map((student) => (
                    <TableRow key={student.id}>
                      <TableCell>{student.full_name}</TableCell>
                      <TableCell>{student.class_name}</TableCell>
                      <TableCell>
                        {student.subject_averages.map((subject, index) => (
                          <div key={index}>
                            {subject.subject_name} (Четверть {subject.quarter}): {subject.average.toFixed(2)}
                          </div>
                        ))}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Paper>
        </Grid>
      </Grid>
    </Container>
  );
}

export default Analytics; 