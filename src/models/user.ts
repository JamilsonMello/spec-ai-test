import { v4 as uuidv4 } from 'uuid';

class User {
  id: string;
  email: string;
  hashed_password: string;

  constructor(email: string, hashed_password: string) {
    this.id = uuidv4();
    this.email = email;
    this.hashed_password = hashed_password;
  }
}

export default User;