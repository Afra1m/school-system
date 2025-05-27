import React, { useState, useEffect } from 'react';
import {
  Container,
  Grid,
  Paper,
  Typography,
  Box,
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

function Dashboard() {
  const [stats, setStats] = useState({
    studentsCount: 0,
    teachersCount: 0,
    averageGrade: 0,
  });
  const [classPerformance, setClassPerformance] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const token = localStorage.getItem('token');
        const headers = {
          'Authorization': `Bearer ${token}`
        };

        const [studentsRes, teachersRes, gradesRes, performanceRes] = await Promise.all([
          axios.get('http://localhost:8000/stats/students-count', { headers }),
          axios.get('http://localhost:8000/stats/teachers-count', { headers }),
          axios.get('http://localhost:8000/stats/average-grade', { headers }),
          axios.get('http://localhost:8000/stats/class-performance', { headers }),
        ]);

        setStats({
          studentsCount: studentsRes.data.count,
          teachersCount: teachersRes.data.count,
          averageGrade: gradesRes.data.average,
        });
        setClassPerformance(performanceRes.data);
        setLoading(false);
      } catch (err) {
        setError('Ошибка при загрузке статистики');
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (loading) return <Typography>Загрузка...</Typography>;
  if (error) return <Typography color="error">{error}</Typography>;

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Grid container spacing={3}>
        {/* Заголовок */}
        <Grid item xs={12}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
            <Typography component="h1" variant="h4" color="primary" gutterBottom>
              Панель управления
            </Typography>
          </Paper>
        </Grid>

        {/* Статистика */}
        <Grid item xs={12} md={4}>
          <Paper
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              height: 240,
            }}
          >
            <Typography component="h2" variant="h6" color="primary" gutterBottom>
              Количество учеников
            </Typography>
            <Typography component="p" variant="h4">
              {stats.studentsCount}
            </Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              height: 240,
            }}
          >
            <Typography component="h2" variant="h6" color="primary" gutterBottom>
              Количество учителей
            </Typography>
            <Typography component="p" variant="h4">
              {stats.teachersCount}
            </Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              height: 240,
            }}
          >
            <Typography component="h2" variant="h6" color="primary" gutterBottom>
              Средний балл
            </Typography>
            <Typography component="p" variant="h4">
              {stats.averageGrade.toFixed(2)}
            </Typography>
          </Paper>
        </Grid>

        {/* График успеваемости */}
        <Grid item xs={12}>
          <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
            <Typography component="h2" variant="h6" color="primary" gutterBottom>
              Успеваемость по классам
            </Typography>
            <Box sx={{ height: 400 }}>
              <ResponsiveContainer width="100%" height="100%">
                <BarChart
                  data={classPerformance}
                  margin={{
                    top: 20,
                    right: 30,
                    left: 20,
                    bottom: 5,
                  }}
                >
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis domain={[0, 5]} />
                  <Tooltip />
                  <Legend />
                  <Bar dataKey="value" name="Средний балл" fill="#8884d8" />
                </BarChart>
              </ResponsiveContainer>
            </Box>
          </Paper>
        </Grid>
      </Grid>
    </Container>
  );
}

export default Dashboard; 