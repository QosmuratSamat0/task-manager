const Task = require('../models/Task');
const Project = require('../models/Project');
const User = require('../models/User');
const { toTaskResponse } = require('../utils/formatters');

async function getTasks(req, res) {
  const query = {};
  
  // If user is authenticated and not admin, restrict to their own tasks
  if (req.user && req.user.role !== 'admin') {
    query.user = req.user._id;
  } else if (req.query.projectId || req.query.project_id) {
    query.project = req.query.projectId || req.query.project_id;
  }
  
  if (req.query.userId || req.query.user_id) {
    query.user = req.query.userId || req.query.user_id;
  }
  
  const tasks = await Task.find(query).sort({ createdAt: -1 });
  return res.json({ data: tasks.map(toTaskResponse) });
}

async function getTaskById(req, res) {
  const task = await Task.findById(req.params.id);
  if (!task) {
    return res.status(404).json({ error: 'Task not found' });
  }
  return res.json({ data: toTaskResponse(task) });
}

async function createTask(req, res) {
  const projectId = req.body.project_id || req.body.project || null;
  const userId = req.body.user_id || req.body.userId || null;

  if (!userId) {
    return res.status(400).json({ error: 'user_id is required' });
  }

  const user = await User.findById(userId);
  if (!user) {
    return res.status(400).json({ error: 'Invalid user id' });
  }

  if (!projectId) {
    // Project is optional for legacy frontend
  } else {
    const project = await Project.findById(projectId);
    if (!project) {
      return res.status(400).json({ error: 'Invalid project id' });
    }
  }

  const task = await Task.create({
    title: req.body.title,
    description: req.body.description || '',
    status: req.body.status || 'todo',
    priority: req.body.priority || 'medium',
    deadline: req.body.deadline,
    user: userId,
    project: projectId || undefined,
  });

  return res.status(201).json({ data: toTaskResponse(task) });
}

async function updateTask(req, res) {
  const updates = {};
  const projectId = req.body.project_id || req.body.project || null;
  const userId = req.body.user_id || req.body.userId || null;

  if (projectId) {
    const project = await Project.findById(projectId);
    if (!project) {
      return res.status(400).json({ error: 'Invalid project id' });
    }
    updates.project = projectId;
  }

  if (userId) {
    const user = await User.findById(userId);
    if (!user) {
      return res.status(400).json({ error: 'Invalid user id' });
    }
    updates.user = userId;
  }

  if (req.body.title !== undefined) updates.title = req.body.title;
  if (req.body.description !== undefined) updates.description = req.body.description;
  if (req.body.status !== undefined) updates.status = req.body.status;
  if (req.body.priority !== undefined) updates.priority = req.body.priority;
  if (req.body.deadline !== undefined) updates.deadline = req.body.deadline;

  if (Object.keys(updates).length === 0) {
    return res.status(400).json({ error: 'No valid fields to update' });
  }

  const task = await Task.findByIdAndUpdate(req.params.id, updates, {
    new: true,
    runValidators: true,
  });
  if (!task) {
    return res.status(404).json({ error: 'Task not found' });
  }
  return res.json({ data: toTaskResponse(task) });
}

async function deleteTask(req, res) {
  const task = await Task.findByIdAndDelete(req.params.id);
  if (!task) {
    return res.status(404).json({ error: 'Task not found' });
  }
  return res.json({ success: true });
}

async function getTasksByUser(req, res) {
  const userId = req.params.userId;
  const tasks = await Task.find({ user: userId }).sort({ createdAt: -1 });
  return res.json({ data: tasks.map(toTaskResponse) });
}

module.exports = {
  getTasks,
  getTaskById,
  createTask,
  updateTask,
  deleteTask,
  getTasksByUser,
};
