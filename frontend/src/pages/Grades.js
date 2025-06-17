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
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Autocomplete,
} from '@mui/material';
import axios from 'axios';

function Grades() {
  const [grades, setGrades] = useState([]);
  const [filteredGrades, setFilteredGrades] = useState([]);
  const [students, setStudents] = useState([]);
  const [subjects, setSubjects] = useState([]);
  const [open, setOpen] = useState(false);
  const [selectedStudent, setSelectedStudent] = useState(null);
  const [selectedSubject, setSelectedSubject] = useState(null);
  const [newGrade, setNewGrade] = useState({
    student_id: '',
    subject_id: '',
    grade: '',
    quarter: '',
  });
  const [editOpen, setEditOpen] = useState(false);
  const [editingGrade, setEditingGrade] = useState(null);

  useEffect(() => {
    fetchGrades();
    fetchStudents();
    fetchSubjects();
  }, []);

  useEffect(() => {
    filterGrades();
  }, [grades, selectedStudent, selectedSubject]);

  const filterGrades = () => {
    let filtered = grades;
    
    if (selectedStudent) {
      filtered = filtered.filter(grade => grade.student_id === selectedStudent.id);
    }
    
    if (selectedSubject) {
      filtered = filtered.filter(grade => grade.subject_id === selectedSubject.id);
    }
    
    setFilteredGrades(filtered);
  };

  const fetchGrades = async () => {
    try {
      const response = await axios.get('http://localhost:8000/grades', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });
      // Сортируем оценки по ID ученика и предмета
      const sortedGrades = response.data.sort((a, b) => {
        if (a.student_id !== b.student_id) {
          return a.student_id - b.student_id;
        }
        return a.subject_id - b.subject_id;
      });
      setGrades(sortedGrades);
    } catch (error) {
      console.error('Ошибка при получении оценок:', error);
    }
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
      console.error('Ошибка при получении списка студентов:', error);
    }
  };

  const fetchSubjects = async () => {
    try {
      const response = await axios.get('http://localhost:8000/subjects', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });
      console.log('Полученные предметы с сервера:', response.data);
      
      const uniqueSubjects = response.data.reduce((acc, current) => {
        const x = acc.find(item => item.id === current.id);
        if (!x) {
          return acc.concat([current]);
        } else {
          return acc;
        }
      }, []);
      console.log('Уникальные предметы после обработки:', uniqueSubjects);
      
      setSubjects(uniqueSubjects);
    } catch (error) {
      console.error('Ошибка при получении списка предметов:', error);
    }
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setNewGrade({
      student_id: '',
      subject_id: '',
      grade: '',
      quarter: '',
    });
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setNewGrade((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await axios.post('http://localhost:8000/grades', newGrade, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });
      handleClose();
      fetchGrades();
    } catch (error) {
      console.error('Ошибка при создании оценки:', error);
    }
  };

  const handleDelete = async (id) => {
    if (!id) {
      console.error('ID оценки не указан');
      return;
    }

    if (window.confirm('Вы уверены, что хотите удалить эту оценку?')) {
      try {
        await axios.delete(`http://localhost:8000/grades/${id}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        });
        console.log(`Оценка с ID ${id} успешно удалена`);
        fetchGrades();
      } catch (error) {
        console.error('Ошибка при удалении оценки:', error);
      }
    }
  };

  const getStudentName = (studentId) => {
    const student = students.find((s) => s.id === studentId);
    return student ? student.full_name : '';
  };

  const getSubjectName = (subjectId) => {
    const subject = subjects.find((s) => s.id === subjectId);
    return subject ? subject.name : '';
  };

  const handleEditClick = (grade) => {
    setEditingGrade(grade);
    setEditOpen(true);
  };

  const handleEditClose = () => {
    setEditOpen(false);
    setEditingGrade(null);
  };

  const handleEditChange = (e) => {
    const { name, value } = e.target;
    setEditingGrade(prev => ({
      ...prev,
      [name]: name === 'grade' || name === 'quarter' ? parseInt(value) || '' : value,
    }));
  };

  const handleEditSubmit = async () => {
    try {
      await axios.put(`http://localhost:8000/grades/${editingGrade.id}`, editingGrade, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });
      handleEditClose();
      fetchGrades();
    } catch (error) {
      console.error('Ошибка при обновлении оценки:', error);
    }
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
        <Typography component="h1" variant="h4" color="primary">
          Управление оценками
        </Typography>
        <Button variant="contained" color="primary" onClick={handleClickOpen}>
          Добавить оценку
        </Button>
      </Box>

      <Box sx={{ mb: 2, display: 'flex', gap: 2 }}>
        <Autocomplete
          options={students}
          getOptionLabel={(option) => `${option.full_name} (${option.class_name})`}
          value={selectedStudent}
          onChange={(event, newValue) => {
            setSelectedStudent(newValue);
          }}
          sx={{ minWidth: 350, maxWidth: 500, flex: 1 }}
          renderInput={(params) => (
            <TextField
              {...params}
              label="Фильтр по ученику"
              variant="outlined"
              fullWidth
            />
          )}
        />
        
        <Autocomplete
          options={subjects}
          getOptionLabel={(option) => option.name}
          isOptionEqualToValue={(option, value) => option.id === value.id}
          value={selectedSubject}
          onChange={(event, newValue) => {
            setSelectedSubject(newValue);
          }}
          sx={{ minWidth: 220, maxWidth: 350 }}
          renderInput={(params) => (
            <TextField
              {...params}
              label="Фильтр по предмету"
              variant="outlined"
              fullWidth
            />
          )}
        />
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Студент</TableCell>
              <TableCell>Предмет</TableCell>
              <TableCell>Оценка</TableCell>
              <TableCell>Четверть</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredGrades.map((grade) => (
              <TableRow key={grade.id}>
                <TableCell>{grade.id}</TableCell>
                <TableCell>{getStudentName(grade.student_id)}</TableCell>
                <TableCell>{getSubjectName(grade.subject_id)}</TableCell>
                <TableCell>{grade.grade}</TableCell>
                <TableCell>{grade.quarter}</TableCell>
                <TableCell>
                  <Button
                    variant="outlined"
                    color="primary"
                    size="small"
                    sx={{ mr: 1 }}
                    onClick={() => handleEditClick(grade)}
                  >
                    Редактировать
                  </Button>
                  <Button
                    variant="outlined"
                    color="error"
                    size="small"
                    onClick={() => handleDelete(grade.id)}
                  >
                    Удалить
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Добавить новую оценку</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
            <Autocomplete
              options={students}
              getOptionLabel={(option) => `${option.full_name} (${option.class_name})`}
              onChange={(event, newValue) => {
                setNewGrade((prev) => ({
                  ...prev,
                  student_id: newValue ? newValue.id : '',
                }));
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Ученик"
                  variant="outlined"
                  required
                />
              )}
            />
            
            <Autocomplete
              options={subjects}
              getOptionLabel={(option) => option.name}
              isOptionEqualToValue={(option, value) => option.id === value.id}
              onChange={(event, newValue) => {
                setNewGrade((prev) => ({
                  ...prev,
                  subject_id: newValue ? newValue.id : '',
                }));
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Предмет"
                  variant="outlined"
                  required
                />
              )}
            />
            
            <TextField
              name="grade"
              label="Оценка"
              type="number"
              value={newGrade.grade}
              onChange={handleChange}
              required
              inputProps={{ 
                min: 2, 
                max: 5,
                step: 1
              }}
              error={Boolean(newGrade.grade && (newGrade.grade < 2 || newGrade.grade > 5))}
              helperText={newGrade.grade && (newGrade.grade < 2 || newGrade.grade > 5) ? "Оценка должна быть от 2 до 5" : ""}
            />
            
            <TextField
              name="quarter"
              label="Четверть"
              type="number"
              value={newGrade.quarter}
              onChange={handleChange}
              required
              inputProps={{ 
                min: 1, 
                max: 4,
                step: 1
              }}
              error={Boolean(newGrade.quarter && (newGrade.quarter < 1 || newGrade.quarter > 4))}
              helperText={newGrade.quarter && (newGrade.quarter < 1 || newGrade.quarter > 4) ? "Четверть должна быть от 1 до 4" : ""}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Отмена</Button>
          <Button onClick={handleSubmit} variant="contained" color="primary">
            Добавить
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog open={editOpen} onClose={handleEditClose}>
        <DialogTitle>Редактировать оценку</DialogTitle>
        <DialogContent>
          <FormControl fullWidth margin="dense">
            <InputLabel>Ученик</InputLabel>
            <Select
              name="student_id"
              value={editingGrade?.student_id || ''}
              label="Ученик"
              onChange={handleEditChange}
            >
              {students.map((student) => (
                <MenuItem key={student.id} value={student.id}>
                  {student.full_name} ({student.class_name})
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <FormControl fullWidth margin="dense">
            <InputLabel>Предмет</InputLabel>
            <Select
              name="subject_id"
              value={editingGrade?.subject_id || ''}
              label="Предмет"
              onChange={handleEditChange}
            >
              {subjects.map((subject) => (
                <MenuItem key={subject.id} value={subject.id}>
                  {subject.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <TextField
            margin="dense"
            name="grade"
            label="Оценка"
            type="number"
            fullWidth
            variant="outlined"
            value={editingGrade?.grade || ''}
            onChange={handleEditChange}
            inputProps={{ 
              min: 2, 
              max: 5,
              step: 1
            }}
            error={Boolean(editingGrade?.grade && (editingGrade.grade < 2 || editingGrade.grade > 5))}
            helperText={editingGrade?.grade && (editingGrade.grade < 2 || editingGrade.grade > 5) ? "Оценка должна быть от 2 до 5" : ""}
          />
          <TextField
            margin="dense"
            name="quarter"
            label="Четверть"
            type="number"
            fullWidth
            variant="outlined"
            value={editingGrade?.quarter || ''}
            onChange={handleEditChange}
            inputProps={{ 
              min: 1, 
              max: 4,
              step: 1
            }}
            error={Boolean(editingGrade?.quarter && (editingGrade.quarter < 1 || editingGrade.quarter > 4))}
            helperText={editingGrade?.quarter && (editingGrade.quarter < 1 || editingGrade.quarter > 4) ? "Четверть должна быть от 1 до 4" : ""}
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

export default Grades; 