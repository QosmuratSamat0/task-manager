const jwt = require('jsonwebtoken');
const User = require('../models/User');
const { toUserResponse } = require('../utils/formatters');

function createToken(user) {
  const secret = process.env.JWT_SECRET;
  if (!secret) {
    throw new Error('JWT_SECRET is required');
  }
  return jwt.sign(
    { id: user._id, role: user.role },
    secret,
    { expiresIn: process.env.JWT_EXPIRES_IN || '1d' }
  );
}

async function register(req, res) {
  const { user_name: userNameInput, email, password, role } = req.body;
  if (!userNameInput || !email || !password) {
    return res.status(400).json({ error: 'Username, email, and password are required' });
  }
  if (role && !['user', 'admin'].includes(role)) {
    return res.status(400).json({ error: 'Invalid role' });
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
    role: role || 'user',
  });

  const token = createToken(user);
  return res.status(201).json({
    data: toUserResponse(user),
    token,
  });
}

async function login(req, res) {
  const { email, user_name: userNameInput, password } = req.body;
  if ((!email && !userNameInput) || !password) {
    return res.status(400).json({ error: 'Username/email and password are required' });
  }

  const user = userNameInput
    ? await User.findOne({ userName: userNameInput.trim() })
    : await User.findOne({ email: email.toLowerCase().trim() });
  if (!user) {
    return res.status(401).json({ error: 'Invalid credentials' });
  }

  const isMatch = await user.comparePassword(password);
  if (!isMatch) {
    return res.status(401).json({ error: 'Invalid credentials' });
  }

  const token = createToken(user);
  return res.json({
    data: toUserResponse(user),
    token,
  });
}

module.exports = { register, login };
