export interface Email {
  messageId: string;
  date: Date;
  from: string;
  to: string[];
  cc: string[];
  bcc: string[];
  subject?: string;
  mimeVersion?: string;
  contentType?: string;
  contentTransferEncoding?: string;
  xFrom?: string;
  xTo?: string[];
  xCc?: string[];
  xBcc?: string[];
  xFolder?: string;
  xOrigin?: string;
  xFileName?: string;
  body: string;
}
