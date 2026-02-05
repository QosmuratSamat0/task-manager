const Category = require('../models/Category');
const { toCategoryResponse } = require('../utils/formatters');
const asyncHandler = require('../utils/asyncHandler');

// Get all categories
const listCategories = asyncHandler(async (req, res) => {
  const categories = await Category.find();
  res.json({ status: 'success', data: categories.map(toCategoryResponse) });
});

// Get category by ID
const getCategory = asyncHandler(async (req, res) => {
  const category = await Category.findById(req.params.id);
  if (!category) return res.status(404).json({ error: 'Category not found' });
  res.json({ status: 'success', data: toCategoryResponse(category) });
});

// Create category (ADMIN only)
const createCategory = asyncHandler(async (req, res) => {
  if (req.user.role !== 'admin') {
    return res.status(403).json({ error: 'Admin access required' });
  }
  const { name, description, color } = req.body;
  const category = new Category({ name, description, color });
  await category.save();
  res.status(201).json({ status: 'success', data: toCategoryResponse(category) });
});

// Update category (ADMIN only)
const updateCategory = asyncHandler(async (req, res) => {
  if (req.user.role !== 'admin') {
    return res.status(403).json({ error: 'Admin access required' });
  }
  const { name, description, color } = req.body;
  const category = await Category.findByIdAndUpdate(
    req.params.id,
    { name, description, color },
    { new: true }
  );
  if (!category) return res.status(404).json({ error: 'Category not found' });
  res.json({ status: 'success', data: toCategoryResponse(category) });
});

// Delete category (ADMIN only)
const deleteCategory = asyncHandler(async (req, res) => {
  if (req.user.role !== 'admin') {
    return res.status(403).json({ error: 'Admin access required' });
  }
  const category = await Category.findByIdAndDelete(req.params.id);
  if (!category) return res.status(404).json({ error: 'Category not found' });
  res.json({ status: 'success', message: 'Category deleted' });
});

module.exports = { listCategories, getCategory, createCategory, updateCategory, deleteCategory };
