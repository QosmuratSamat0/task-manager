require('dotenv').config();
const createApp = require('./src/app');
const { connectDB } = require('./src/config/db');
const User = require('./src/models/User');

const PORT = process.env.PORT || 3000;
const MONGODB_URI = process.env.MONGODB_URI;

const app = createApp();

async function initializeAdmin() {
  try {
    // Check if admin user exists
    const adminExists = await User.findOne({ role: 'admin' });
    if (!adminExists) {
      // Create default admin user
      const admin = await User.create({
        userName: 'admin',
        email: 'admin@taskmanager.local',
        password: 'secret123',
        role: 'admin',
      });
      console.log('Admin user created: admin@taskmanager.local / secret123');
    }
  } catch (err) {
    console.error('Failed to initialize admin user', err);
  }
}

connectDB(MONGODB_URI)
  .then(async () => {
    console.log('Connected to MongoDB');
    await initializeAdmin();
    app.listen(PORT, () => {
      console.log(`Server is running on port ${PORT}`);
    });
  })
  .catch((err) => {
    console.error('Failed to connect to MongoDB', err);
    process.exit(1);
  });
