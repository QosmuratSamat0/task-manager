const express = require('express');
const cors = require('cors');

const authRoutes = require('./routes/authRoutes');
const taskRoutes = require('./routes/taskRoutes');
const projectRoutes = require('./routes/projectRoutes');
const userRoutes = require('./routes/userRoutes');
const errorHandler = require('./middleware/errorHandler');

function createApp() {
  const app = express();

  app.use(cors());
  app.use(express.json());

  app.get('/status', (req, res) => {
    return res.status(200).json({ status: 'OK' });
  });

  app.use('/auth', authRoutes);
  app.use('/users', userRoutes);
  app.use('/tasks', taskRoutes);
  app.use('/projects', projectRoutes);

  app.use(errorHandler);

  return app;
}

module.exports = createApp;
