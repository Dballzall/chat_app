import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ChatService, ChatMessage } from '../../services/chat.service';

@Component({
  selector: 'app-chat',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './chat.component.html',
  styleUrl: './chat.component.css'
})
export class ChatComponent implements OnInit, OnDestroy {
  username = '';
  message = '';
  messages: ChatMessage[] = [];
  isConnected = false;

  constructor(private chatService: ChatService) {}

  ngOnInit() {
    this.chatService.getMessages().subscribe(msg => {
      this.messages.push(msg);
    });
  }

  connect() {
    if (this.username.trim()) {
      this.chatService.connect(this.username);
      this.isConnected = true;
    }
  }

  sendMessage() {
    if (this.message.trim()) {
      this.chatService.sendMessage(this.message);
      this.message = '';
    }
  }

  ngOnDestroy() {
    this.chatService.disconnect();
  }
}
