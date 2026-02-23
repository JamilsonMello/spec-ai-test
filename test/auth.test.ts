import request from 'supertest';
import app from '../src/app';
import User from '../src/models/user';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';

describe('Auth Endpoints', () => {
  let user: User;
  let token: string;

  beforeAll(async () => {
    // Create a test user
    const hashedPassword = await bcrypt.hash('password123', 10);
    user = await User.create({ email: 'test@example.com', hashed_password: hashedPassword });

    // Generate a JWT for the test user
    token = jwt.sign({ id: user.id }, process.env.JWT_SECRET || 'secret', { expiresIn: '1h' });
  });

  afterAll(async () => {
    // Delete the test user
    await User.destroy({ where: { id: user.id } });
  });

  describe('POST /auth/login', () => {
    it('should return a JWT for valid credentials', async () => {
      const res = await request(app)
        .post('/auth/login')
        .send({ email: 'test@example.com', password: 'password123' });

      expect(res.statusCode).toEqual(200);
      expect(res.body).toHaveProperty('token');
    });

    it('should return 401 for invalid email', async () => {
      const res = await request(app)
        .post('/auth/login')
        .send({ email: 'invalid@example.com', password: 'password123' });

      expect(res.statusCode).toEqual(401);
    });

    it('should return 401 for invalid password', async () => {
      const res = await request(app)
        .post('/auth/login')
        .send({ email: 'test@example.com', password: 'wrongpassword' });

      expect(res.statusCode).toEqual(401);
    });

    it('should return 400 for invalid email format', async () => {
      const res = await request(app)
        .post('/auth/login')
        .send({ email: 'invalid-email', password: 'password123' });

      expect(res.statusCode).toEqual(400);
    });

    it('should return 401 for missing email or password', async () => {
      const res = await request(app)
        .post('/auth/login')
        .send({});

      expect(res.statusCode).toEqual(401);
    });
  });
});