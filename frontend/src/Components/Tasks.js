import React, { useState, useEffect } from "react";
import { Container, TextField, Button, Table, TableHead, TableRow, TableCell, TableBody, Paper, IconButton } from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete"; 
import Navbar from "../Components/Navbar";

const API_URL = "http://localhost:30080";

const Tasks = () => {
    const [tasks, setTasks] = useState([]);
    const [newTask, setNewTask] = useState("");
    const [deadline, setDeadline] = useState("");
  
    useEffect(() => {
      fetchTasks();
    }, []);
  
    // ✅ Получение списка задач
    const fetchTasks = async () => {
      const response = await fetch(`${API_URL}/tasks`, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
  
      if (response.ok) {
        var data = await response.json();
        if (data != null) {
            setTasks(data);
        } else {
            setTasks([])
        }
      }
    };
  
    // ✅ Метод для добавления новой задачи
    const handleAddTask = async () => {
      if (!newTask || !deadline) {
        alert("Заполните все поля!");
        return;
      }
  
      const response = await fetch(`${API_URL}/tasks`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({ title: newTask, deadline }),
      });
  
      if (response.ok) {
        setTasks([...tasks, newTask]); // Добавляем задачу в список
        setNewTask("");
        setDeadline("");
        fetchTasks();
      } else {
        alert("Ошибка при добавлении задачи!");
      }
    };

    const handleDeleteTask = async (taskTitle, taskDeadline) => {
        console.log(JSON.stringify({ title: taskTitle, deadline: taskDeadline }))
        const response = await fetch(`${API_URL}/tasks`, {
            method: "DELETE",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
            body: JSON.stringify({ title: taskTitle, deadline: formatDate(taskDeadline) }),
            mode: "cors",
          });
    
        if (response.ok) {
            fetchTasks();
        } else {
          alert("Ошибка при удалении задачи!");
        }
      };

    // Функция для форматирования даты (оставляет только YYYY-MM-DD)
  const formatDate = (dateString) => {
    if (dateString != null) {
        return dateString.split("T")[0]; // Берем только дату до "T"
    }
  };
  
    return (
      <Container>
        <Navbar />
        <Paper sx={{ padding: 2, marginTop: 2 }}>
          <TextField label="Новая задача" value={newTask} onChange={(e) => setNewTask(e.target.value)} />
          <TextField type="date" value={deadline} onChange={(e) => setDeadline(e.target.value)} />
          <Button variant="contained" onClick={handleAddTask}>
            Добавить
          </Button>
        </Paper>
  
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Задача</TableCell>
              <TableCell>Дедлайн</TableCell>
              <TableCell></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {tasks.map((task, index) => (
              <TableRow key={index}>
                <TableCell>{task.title}</TableCell>
                <TableCell>{formatDate(task.deadline)}</TableCell>
                <TableCell>
                    <IconButton color="error" onClick={() => handleDeleteTask(task.title, task.deadline)}>
                        <DeleteIcon />
                    </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Container>
    );
  };
  
  export default Tasks;