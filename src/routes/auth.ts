import express, { Request, Response } from 'express';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';
import { User } from '../models/user';

const router = express.Router();

router.post('/login', async (req: Request, res: Response) => {
  const { email, password } = req.body;

  // FR-003: Validate email format
  if (!email || !/^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/.test(email)) {
    return res.status(400).send({ message: 'Invalid email format' });
  }

  try {
    // FR-006: If credentials are valid, generate and return an authentication token
    const user = await User.findOne({ email });

    if (!user) {
      return res
        .status(401)
        .send({ message: 'email ou senha inválidos' });
    }

    // FR-004: Compare the password with stored hash
    const passwordMatch = await bcrypt.compare(password, user.password);

    if (!passwordMatch) {
      return res
        .status(401)
        .send({ message: 'email ou senha inválidos' });
    }

    // FR-005: Generate and return an authentication token
    const token = jwt.sign({ userId: user._id }, 'secret', { expiresIn: '1h' });

    res.status(200).send({ token });
  } catch (error) {
    console.error(error);
    return res.status(500).send({ message: 'Internal server error' });
  }
});

export { router as authRouter };
