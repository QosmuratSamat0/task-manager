const User = require('../models/User');
const Task = require('../models/Task');
const Project = require('../models/Project');
const asyncHandler = require('../utils/asyncHandler');

// Get server stats (ADMIN only)
const getStats = asyncHandler(async (req, res) => {
  if (req.user.role !== 'admin') {
    return res.status(403).json({ error: 'Admin access required' });
  }

  const totalUsers = await User.countDocuments();
  const totalTasks = await Task.countDocuments();
  const totalProjects = await Project.countDocuments();
  
  const tasksByStatus = await Task.aggregate([
    { $group: { _id: '$status', count: { $sum: 1 } } },
  ]);
  
  const tasksByPriority = await Task.aggregate([
    { $group: { _id: '$priority', count: { $sum: 1 } } },
  ]);

  const recentTasks = await Task.find()
    .populate('user', 'userName')
    .populate('project', 'name')
    .sort({ createdAt: -1 })
    .limit(10);

  const recentUsers = await User.find()
    .sort({ createdAt: -1 })
    .limit(10)
    .select('-password');

  res.json({
    status: 'success',
    data: {
      totalUsers,
      totalTasks,
      totalProjects,
      tasksByStatus,
      tasksByPriority,
      recentTasks,
      recentUsers,
      timestamp: new Date(),
    },
  });
});

module.exports = { getStats };
