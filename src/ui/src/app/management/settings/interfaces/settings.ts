export interface ISettings {
  baseUrl: string;
  smtp: {
    host: string;
    port: number;
    username: string;
    password: string;
  };
  sender: {
    enable: boolean;
    email: string;
    orderRecipients: string[];
  };
}
