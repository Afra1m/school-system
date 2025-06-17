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
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import axios from 'axios';

function Students() {
  const [students, setStudents] = useState([]);
  const [filteredStudents, setFilteredStudents] = useState([]);
  const [open, setOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [editingStudent, setEditingStudent] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedClass, setSelectedClass] = useState('');
  const [newStudent, setNewStudent] = useState({
    full_name: '',
    class_name: '',
  });

  useEffect(() => {
    fetchStudents();
  }, []);

  useEffect(() => {
    filterStudents();
  }, [students, searchQuery, selectedClass]);

  const filterStudents = () => {
    let filtered = [...students];
    
    // Фильтр по поисковому запросу
    if (searchQuery) {
      filtered = filtered.filter(student => 
        student.full_name.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }
    
    // Фильтр по классу
    if (selectedClass) {
      filtered = filtered.filter(student => 
        student.class_name === selectedClass
      );
    }
    
    setFilteredStudents(filtered);
  };

  const fetchStudents = async () => {
    try {
      const response = await axios.get('http://localhost:8000/students', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });
      setStudents(response.data);
    } catch (error) {
      console.error('Ошибка при получении списка учеников:', error);
    }
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setNewStudent({ full_name: '', class_name: '' });
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setNewStudent((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async () => {
    try {
      await axios.post(
        'http://localhost:8000/students',
        newStudent,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        }
      );
      handleClose();
      fetchStudents();
    } catch (error) {
      console.error('Ошибка при создании ученика:', error);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить этого ученика?')) {
      try {
        await axios.delete(`http://localhost:8000/students/${id}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        });
        console.log(`Ученик с ID ${id} успешно удален`);
        fetchStudents();
      } catch (error) {
        console.error('Ошибка при удалении ученика:', error);
      }
    }
  };

  const handleEditClick = (student) => {
    setEditingStudent(student);
    setEditOpen(true);
  };

  const handleEditClose = () => {
    setEditOpen(false);
    setEditingStudent(null);
  };

  const handleEditChange = (e) => {
    const { name, value } = e.target;
    setEditingStudent(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleEditSubmit = async () => {
    try {
      await axios.put(`http://localhost:8000/students/${editingStudent.id}`, editingStudent, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });
      handleEditClose();
      fetchStudents();
    } catch (error) {
      console.error('Ошибка при обновлении ученика:', error);
    }
  };

  // Получаем уникальные классы для фильтра
  const uniqueClasses = [...new Set(students.map(student => student.class_name))].sort();

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
        <Typography component="h1" variant="h4" color="primary">
          Список учеников
        </Typography>
        <Button variant="contained" color="primary" onClick={handleClickOpen}>
          Добавить ученика
        </Button>
      </Box>

      {/* Фильтры */}
      <Box sx={{ mb: 2, display: 'flex', gap: 2 }}>
        <TextField
          label="Поиск по имени"
          variant="outlined"
          size="small"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          sx={{ width: 300 }}
        />
        <FormControl size="small" sx={{ width: 200 }}>
          <InputLabel>Класс</InputLabel>
          <Select
            value={selectedClass}
            label="Класс"
            onChange={(e) => setSelectedClass(e.target.value)}
          >
            <MenuItem value="">Все классы</MenuItem>
            {uniqueClasses.map((className) => (
              <MenuItem key={className} value={className}>
                {className}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>ФИО</TableCell>
              <TableCell>Класс</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredStudents.map((student) => (
              <TableRow key={student.id}>
                <TableCell>{student.id}</TableCell>
                <TableCell>{student.full_name}</TableCell>
                <TableCell>{student.class_name}</TableCell>
                <TableCell>
                  <Button
                    variant="outlined"
                    color="primary"
                    size="small"
                    sx={{ mr: 1 }}
                    onClick={() => handleEditClick(student)}
                  >
                    Редактировать
                  </Button>
                  <Button
                    variant="outlined"
                    color="error"
                    size="small"
                    onClick={() => handleDelete(student.id)}
                  >
                    Удалить
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Диалог добавления ученика */}
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Добавить нового ученика</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            name="full_name"
            label="ФИО"
            type="text"
            fullWidth
            variant="outlined"
            value={newStudent.full_name}
            onChange={handleChange}
          />
          <TextField
            margin="dense"
            name="class_name"
            label="Класс"
            type="text"
            fullWidth
            variant="outlined"
            value={newStudent.class_name}
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

      {/* Диалог редактирования ученика */}
      <Dialog open={editOpen} onClose={handleEditClose}>
        <DialogTitle>Редактировать ученика</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            name="full_name"
            label="ФИО"
            type="text"
            fullWidth
            variant="outlined"
            value={editingStudent?.full_name}
            onChange={handleEditChange}
          />
          <TextField
            margin="dense"
            name="class_name"
            label="Класс"
            type="text"
            fullWidth
            variant="outlined"
            value={editingStudent?.class_name}
            onChange={handleEditChange}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleEditClose}>Отмена</Button>
          <Button onClick={handleEditSubmit} variant="contained" color="primary">
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}

export default Students; 