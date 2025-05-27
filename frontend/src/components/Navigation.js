import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  Container,
} from '@mui/material';

function Navigation({ onLogout }) {
  return (
    <AppBar position="static">
      <Container maxWidth="lg">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Школьная система
          </Typography>
          <Box sx={{ display: 'flex', gap: 2 }}>
            <Button
              color="inherit"
              component={RouterLink}
              to="/"
            >
              Главная
            </Button>
            <Button
              color="inherit"
              component={RouterLink}
              to="/students"
            >
              Ученики
            </Button>
            <Button
              color="inherit"
              component={RouterLink}
              to="/teachers"
            >
              Учителя
            </Button>
            <Button
              color="inherit"
              component={RouterLink}
              to="/grades"
            >
              Оценки
            </Button>
            <Button
              color="inherit"
              component={RouterLink}
              to="/analytics"
            >
              Аналитика
            </Button>
            <Button
              color="inherit"
              onClick={onLogout}
            >
              Выйти
            </Button>
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
}

export default Navigation; 