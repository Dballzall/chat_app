// src/app/services/chat.service.ts

import { Injectable } from '@angular/core';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { Observable, Subject } from 'rxjs';

export interface ChatMessage {
  type: string;
  content: string;
  username: string;
}

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private socket$!: WebSocketSubject<ChatMessage>;
  private messagesSubject = new Subject<ChatMessage>();

  constructor() {}

  connect(username: string): void {
    this.socket$ = webSocket(`ws://localhost:8080/ws?username=${username}`);
    
    this.socket$.subscribe({
      next: (message) => this.messagesSubject.next(message),
      error: (err) => console.error(err),
      complete: () => console.log('WebSocket connection closed')
    });
  }

  sendMessage(content: string): void {
    this.socket$.next({ type: 'message', content, username: '' });
  }

  getMessages(): Observable<ChatMessage> {
    return this.messagesSubject.asObservable();
  }

  disconnect(): void {
    if (this.socket$) {
      this.socket$.complete();
    }
  }
}