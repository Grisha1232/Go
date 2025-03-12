const API_URL = "http://localhost:8080";
let authToken = localStorage.getItem("token") || "";

// Проверка входа при загрузке
document.addEventListener("DOMContentLoaded", () => {
    if (authToken) {
        document.getElementById("auth-section").style.display = "none";
        document.getElementById("tasks-section").style.display = "block";
        fetchTasks();
    }
});

async function register() {
    const username = document.getElementById("register-username").value;
    const password = document.getElementById("register-password").value;

    const response = await fetch(`${API_URL}/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });

    if (response.ok) {
        document.getElementById("register-message").textContent = "Регистрация успешна!";
    } else {
        document.getElementById("register-message").textContent = "Ошибка регистрации!";
    }
}

async function login() {
    const username = document.getElementById("login-username").value;
    const password = document.getElementById("login-password").value;

    const response = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });

    if (response.ok) {
        const data = await response.json();
        authToken = data.token;
        localStorage.setItem("token", authToken);
        document.getElementById("auth-section").style.display = "none";
        document.getElementById("tasks-section").style.display = "block";
        fetchTasks();
    } else {
        document.getElementById("auth-message").textContent = "Ошибка входа!";
    }
}

function logout() {
    localStorage.removeItem("token");
    authToken = "";
    document.getElementById("auth-section").style.display = "block";
    document.getElementById("tasks-section").style.display = "none";
}

async function fetchTasks() {
    const response = await fetch(`${API_URL}/tasks`, {
        method: "GET",
        headers: { "Authorization": `Bearer ${authToken}` }
    });

    if (response.ok) {
        const tasks = await response.json();
        const taskList = document.getElementById("task-list");
        taskList.innerHTML = "";
        tasks.forEach(task => {
            const li = document.createElement("li");
            li.innerHTML = `
                <strong>${task.title}</strong>: ${task.description} (до ${new Date(task.deadline).toLocaleString()}) 
                <button onclick="deleteTask(${task.id})">Удалить</button>
            `;
            taskList.appendChild(li);
        });
    } else {
        console.error("Ошибка загрузки задач");
    }
}

async function createTask() {
    const title = document.getElementById("task-title").value;
    const description = document.getElementById("task-desc").value;
    const deadline = document.getElementById("task-deadline").value;

    const response = await fetch(`${API_URL}/tasks`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${authToken}`
        },
        body: JSON.stringify({ title, description, deadline, completed: false })
    });

    if (response.ok) {
        document.getElementById("task-message").textContent = "Задача добавлена!";
        fetchTasks();
    } else {
        document.getElementById("task-message").textContent = "Ошибка добавления задачи.";
    }
}
