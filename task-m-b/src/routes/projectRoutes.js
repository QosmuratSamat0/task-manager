const express = require('express');
const asyncHandler = require('../utils/asyncHandler');
const projectController = require('../controllers/projectController');

const router = express.Router();

router.get('/', asyncHandler(projectController.getProjects));
router.get('/:id', asyncHandler(projectController.getProjectById));
router.post('/', asyncHandler(projectController.createProject));
router.put('/:id', asyncHandler(projectController.updateProject));
router.delete('/:id', asyncHandler(projectController.deleteProject));

module.exports = router;
