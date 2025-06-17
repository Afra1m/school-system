import React, { useState, useEffect } from 'react';
import {
  Paper,
  Typography,
  Box,
  CircularProgress,
} from '@mui/material';
import axios from 'axios';

function TeacherPerformance() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://localhost:8000/stats/teacher-performance', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        });
        setData(response.data);
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
        Успеваемость учителей
      </Typography>
      <Box sx={{ mt: 2 }}>
        <Typography variant="subtitle1" color="primary">
          Лучший учитель:
        </Typography>
        <Typography>
          {data?.best_teacher?.teacher_name} - {data?.best_teacher?.average_grade}
        </Typography>
      </Box>
      <Box sx={{ mt: 2 }}>
        <Typography variant="subtitle1" color="error">
          Худший учитель:
        </Typography>
        <Typography>
          {data?.worst_teacher?.teacher_name} - {data?.worst_teacher?.average_grade}
        </Typography>
      </Box>
    </Paper>
  );
}

export default TeacherPerformance; 