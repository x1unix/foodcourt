export interface ISettings {
  baseUrl: string;
  smtp: {
    host: string;
    port: number;
    username: string;
    password: string;
  };
  sender: {
    email: string;
    orderRecipients: string[];
  };
}
