import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { TextField, Button, Container, Paper, Typography, Box } from "@mui/material";

const API_URL = "http://localhost:30080";

const Login = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      const response = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        alert(`Ошибка входа: ${errorText}`);
        return;
      }

      const data = await response.json();
      console.log("Вход " + username + " " + data.token)
      localStorage.setItem("token", data.token);
      localStorage.setItem("username", username);
      navigate("/tasks");
    } catch (error) {
      alert("Ошибка сети: " + error.message);
    }
  };

  return (
    <Container maxWidth="xs">
      <Paper elevation={3} sx={{ padding: 4, textAlign: "center", marginTop: "100px" }}>
        <Typography variant="h5">Вход</Typography>
        <form onSubmit={handleLogin}>
          <TextField fullWidth label="Логин" margin="normal" value={username} onChange={(e) => setUsername(e.target.value)} required />
          <TextField fullWidth label="Пароль" type="password" margin="normal" value={password} onChange={(e) => setPassword(e.target.value)} required />
          <Button fullWidth variant="contained" color="primary" type="submit" sx={{ marginTop: 2 }}>
            Войти
          </Button>
        </form>
        <Box sx={{ marginTop: 2 }}>
          <Button onClick={() => navigate("/register")}>Нет аккаунта? Регистрация</Button>
        </Box>
      </Paper>
    </Container>
  );
};

export default Login;
