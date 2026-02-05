const express = require('express');
const { listCategories, getCategory, createCategory, updateCategory, deleteCategory } = require('../controllers/categoryController');
const { authenticate } = require('../middleware/auth');

const router = express.Router();

router.get('/', listCategories);
router.get('/:id', getCategory);
router.post('/', authenticate, createCategory);
router.put('/:id', authenticate, updateCategory);
router.delete('/:id', authenticate, deleteCategory);

module.exports = router;
