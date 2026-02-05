const mongoose = require('mongoose');

const categorySchema = new mongoose.Schema(
  {
    name: { type: String, required: true, trim: true, maxlength: 120 },
    description: { type: String, trim: true, maxlength: 500 },
    color: { type: String, default: '#007bff' },
  },
  { timestamps: true }
);

module.exports = mongoose.model('Category', categorySchema);
