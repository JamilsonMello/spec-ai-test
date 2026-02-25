import mongoose, { Schema, Document } from 'mongoose';
import { DomainException } from '../exceptions/DomainException';

export interface UserDocument extends Document {
  email: string;
  password?: string;
  name: string;
  role: string;
}

const UserSchema: Schema = new Schema(
  {
    email: {
      type: String,
      required: true,
      unique: true,
      lowercase: true,
      trim: true,
      validate: {
        validator: (email: string) => /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/.test(email),
        message: (props: any) => `${props.value} is not a valid email format!`,
      },
    },
    password: {
      type: String,
      required: true,
      validate: {
        validator: (password: string) => password.length >= 8,
        message: 'Password must be at least 8 characters long',
      },
    },
    name: {
      type: String,
      required: true,
      trim: true,
      validate: {
        validator: (name: string) => name.length >= 3,
        message: 'Name must be at least 3 characters long',
      },
    },
    role: {
      type: String,
      required: true,
      default: 'user',
      enum: ['user', 'admin'],
    },
    passwordRecoveryToken: {
      type: String,
      default: null,
    },
  },
  {
    timestamps: true,
  }
);

export const User = mongoose.model<UserDocument>('User', UserSchema);