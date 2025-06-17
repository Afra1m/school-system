import React, { useState, useEffect } from 'react';
import {
  Paper,
  Typography,
  Box,
  CircularProgress,
  Divider,
} from '@mui/material';
import axios from 'axios';

function ClassPerformance() {
  const [classData, setClassData] = useState(null);
  const [teacherData, setTeacherData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [classResponse, teacherResponse] = await Promise.all([
          axios.get('http://localhost:8000/stats/top-worst-classes', {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
          }),
          axios.get('http://localhost:8000/stats/teacher-performance', {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
          }),
        ]);
        setClassData(classResponse.data);
        setTeacherData(teacherResponse.data);
        setLoading(false);
      } catch (err) {
        console.error('Ошибка при получении данных:', err);
        setError('Не удалось загрузить данные');
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Paper sx={{ p: 2, textAlign: 'center' }}>
        <Typography color="error">{error}</Typography>
      </Paper>
    );
  }

  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Лучшие и худшие показатели
      </Typography>
      
      {/* Классы */}
      <Box sx={{ mt: 2 }}>
        <Typography variant="subtitle1" color="primary" gutterBottom>
          Классы:
        </Typography>
        <Typography>
          Лучший класс: {classData?.top_class}
        </Typography>
        <Typography>
          Худший класс: {classData?.worst_class}
        </Typography>
      </Box>

      <Divider sx={{ my: 2 }} />

      {/* Учителя */}
      <Box sx={{ mt: 2 }}>
        <Typography variant="subtitle1" color="primary" gutterBottom>
          Учителя:
        </Typography>
        <Typography>
          Лучший учитель: {teacherData?.best_teacher?.teacher_name} - {teacherData?.best_teacher?.average_grade}
        </Typography>
        <Typography>
          Худший учитель: {teacherData?.worst_teacher?.teacher_name} - {teacherData?.worst_teacher?.average_grade}
        </Typography>
      </Box>
    </Paper>
  );
}

export default ClassPerformance; 