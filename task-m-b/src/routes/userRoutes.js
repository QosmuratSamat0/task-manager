const express = require('express');
const asyncHandler = require('../utils/asyncHandler');
const userController = require('../controllers/userController');

const router = express.Router();

router.get('/all', asyncHandler(userController.listUsers));
router.get('/:userName', asyncHandler(userController.getUserByName));
router.post('/', asyncHandler(userController.createUser));
router.delete('/:userName', asyncHandler(userController.deleteUser));

module.exports = router;
