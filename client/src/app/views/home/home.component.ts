import { Component, OnInit } from '@angular/core';
import { Event, EventType } from 'src/app/core/model/Chat';
import { ChatService } from 'src/app/core/ws/chat/chat.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  userName = '';
  hasSubmittedUsername = false;

  events : Event[] = [];

  message : Event = {
    type : EventType.MessageEvent,
    data :{
      client : this.userName,
      message : '',
      timestamp : null
    }  
  }
  
  onDisconnect(){
    this.chatService.close();
    this.hasSubmittedUsername = false;
    this.userName = '';
  }

  constructor(private chatService : ChatService) { }

  ngOnInit(): void {
  }

  onConnect(){
    this.chatService.confConnection(this.userName);
    this.hasSubmittedUsername = true;
    this.chatService.getWS()?.subscribe({
      next : (event : Event) => {
        this.events.push(event);
        console.log(this.events);
        console.log(event);
      },
      error : (err) => {console.log(err)},
      complete : () => {console.log('complete')}
    });
  }

  onSend(){
    this.chatService.send(this.message);
    this.message.data.message = '';
  }

}
