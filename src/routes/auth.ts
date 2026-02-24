import express, { Request, Response } from 'express';
import { loginHandler } from '../../gateway/handler/login_handler';
import { registerHandler } from '../../gateway/handler/register_handler';

const router = express.Router();

router.post('/login', loginHandler);
router.post('/register', registerHandler);

export { router as authRouter };
