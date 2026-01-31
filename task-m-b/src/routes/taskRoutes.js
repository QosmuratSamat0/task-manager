const express = require('express');
const asyncHandler = require('../utils/asyncHandler');
const taskController = require('../controllers/taskController');

const router = express.Router();

router.get('/all', asyncHandler(taskController.getTasks));
router.get('/by-user/:userId', asyncHandler(taskController.getTasksByUser));
router.get('/', asyncHandler(taskController.getTasks));
router.get('/:id', asyncHandler(taskController.getTaskById));
router.post('/', asyncHandler(taskController.createTask));
router.put('/:id', asyncHandler(taskController.updateTask));
router.delete('/:id', asyncHandler(taskController.deleteTask));

module.exports = router;
