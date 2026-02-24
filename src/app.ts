import express from 'express';
import mongoose from 'mongoose';
import { authRouter } from './routes/auth';

const app = express();
const port = 3000;

app.use(express.json());

// Connect to MongoDB
mongoose.connect('mongodb://localhost:27017/login-route', {
  useNewUrlParser: true,
  useUnifiedTopology: true,
} as mongoose.ConnectOptions).then(() => {
  console.log('Connected to MongoDB');
}).catch(err => {
  console.error('Could not connect to MongoDB', err);
});

app.use('/api/auth', authRouter);

app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});