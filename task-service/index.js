const express = require("express");
const jwt = require("jsonwebtoken");
const axios = require("axios")

const app = express();
app.use(express.json());

if (!SECRET) {
    console.error("FATAL ERROR: JWT_SECRET is not defined.");
    process.exit(1);
}

let tasks = [];


function authMiddleware(req, res, next) {
  const header = req.headers["authorization"];

  if (!header) return res.status(401).json({ error: "No token" });

  const token = header.split(" ")[1];

  try {
    const decoded = jwt.verify(token, SECRET);
    req.user = decoded.username;
    next();
  } catch (err) {
    return res.status(403).json({ error: "Invalid token" });
  }
}

// Health check
app.get("/", (req, res) => {
  res.json({ message: "Task Service Running" });
});

// Create task (protected)
app.post("/tasks", authMiddleware, async (req, res) => {
  const task = {
    id: tasks.length + 1,
    title: req.body.title,
    user: req.user,
  };

  tasks.push(task);

  try {
    // Ping the Notify Service using the Docker network hostname
    await axios.post('http://notify-service:8080/notify', {
      event: 'task_created',
      taskId: task.id,
      title: task.title,
      user: task.user
    });
    console.log(`[Task Service] Successfully notified event for Task ID: ${task.id}`);
  } catch (error) {
    // If Notify Service is down, we catch the error here so the user still gets their 200 OK response
    console.error(`[Task Service] Failed to reach Notify Service: ${error.message}`);
  }

  res.json(task);
});

// Get tasks (protected)
app.get("/tasks", authMiddleware, (req, res) => {
  const userTasks = tasks.filter(t => t.user === req.user);
  res.json(userTasks);
});

app.listen(3000, () => {
  console.log("Task service running on port 3000");
});
