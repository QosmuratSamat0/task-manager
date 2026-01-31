const User = require('../models/User');
const Task = require('../models/Task');
const { toUserResponse } = require('../utils/formatters');

async function listUsers(req, res) {
  const users = await User.find().sort({ createdAt: -1 });
  return res.json({ data: users.map(toUserResponse) });
}

async function getUserByName(req, res) {
  const userName = req.params.userName;
  const user = await User.findOne({ userName });
  if (!user) {
    return res.status(404).json({ error: 'User not found' });
  }
  return res.json({ data: toUserResponse(user) });
}

async function createUser(req, res) {
  const { user_name: userNameInput, email, password } = req.body;
  if (!userNameInput || !email || !password) {
    return res.status(400).json({ error: 'Username, email, and password are required' });
  }

  const normalizedEmail = email.toLowerCase().trim();
  const normalizedUserName = userNameInput.trim();
  const existing = await User.findOne({
    $or: [{ email: normalizedEmail }, { userName: normalizedUserName }],
  });
  if (existing) {
    return res.status(409).json({ error: 'Email or username already in use' });
  }

  const user = await User.create({
    userName: normalizedUserName,
    email: normalizedEmail,
    password,
    role: 'user',
  });

  return res.status(201).json({
    status: 'created',
    data: toUserResponse(user),
  });
}

async function deleteUser(req, res) {
  const userName = req.params.userName;
  const user = await User.findOneAndDelete({ userName });
  if (!user) {
    return res.status(404).json({ error: 'User not found' });
  }
  await Task.deleteMany({ user: user._id });
  return res.json({ success: true });
}

module.exports = {
  listUsers,
  getUserByName,
  createUser,
  deleteUser,
};
