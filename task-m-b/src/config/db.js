const mongoose = require('mongoose');

async function connectDB(uri) {
  if (!uri) {
    throw new Error('MONGODB_URI is required');
  }
  mongoose.set('strictQuery', true);
  return mongoose.connect(uri);
}

module.exports = { connectDB };
