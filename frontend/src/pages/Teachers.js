import React, { useState, useEffect } from 'react';
import {
  Container,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Box,
} from '@mui/material';
import axios from 'axios';

function Teachers() {
  const [teachers, setTeachers] = useState([]);
  const [open, setOpen] = useState(false);
  const [newTeacher, setNewTeacher] = useState({
    full_name: '',
    room_number: '',
  });

  useEffect(() => {
    fetchTeachers();
  }, []);

  const fetchTeachers = async () => {
    try {
      const response = await axios.get('http://localhost:8000/teachers', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });
      setTeachers(response.data);
    } catch (error) {
      console.error('Ошибка при получении списка учителей:', error);
    }
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setNewTeacher({ full_name: '', room_number: '' });
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setNewTeacher((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await axios.post('http://localhost:8000/teachers', newTeacher, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });
      handleClose();
      fetchTeachers();
    } catch (error) {
      console.error('Ошибка при создании учителя:', error);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить этого учителя?')) {
      try {
        await axios.delete(`http://localhost:8000/teachers/${id}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        });
        console.log(`Учитель с ID ${id} успешно удален`);
        fetchTeachers();
      } catch (error) {
        console.error('Ошибка при удалении учителя:', error);
      }
    }
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
        <Typography component="h1" variant="h4" color="primary">
          Список учителей
        </Typography>
        <Button variant="contained" color="primary" onClick={handleClickOpen}>
          Добавить учителя
        </Button>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>ФИО</TableCell>
              <TableCell>Кабинет</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {teachers.map((teacher) => (
              <TableRow key={teacher.id}>
                <TableCell>{teacher.id}</TableCell>
                <TableCell>{teacher.full_name}</TableCell>
                <TableCell>{teacher.room_number}</TableCell>
                <TableCell>
                  <Button
                    variant="outlined"
                    color="primary"
                    size="small"
                    sx={{ mr: 1 }}
                  >
                    Редактировать
                  </Button>
                  <Button
                    variant="outlined"
                    color="error"
                    size="small"
                    onClick={() => handleDelete(teacher.id)}
                  >
                    Удалить
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Диалог добавления учителя */}
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Добавить нового учителя</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            name="full_name"
            label="ФИО"
            type="text"
            fullWidth
            variant="outlined"
            value={newTeacher.full_name}
            onChange={handleChange}
          />
          <TextField
            margin="dense"
            name="room_number"
            label="Номер кабинета"
            type="text"
            fullWidth
            variant="outlined"
            value={newTeacher.room_number}
            onChange={handleChange}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Отмена</Button>
          <Button onClick={handleSubmit} variant="contained" color="primary">
            Добавить
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}

export default Teachers; 