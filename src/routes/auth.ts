import express, { Request, Response } from 'express';
import { loginHandler } from '../../gateway/handler/login_handler';

const router = express.Router();

router.post('/login', loginHandler);

export { router as authRouter };
