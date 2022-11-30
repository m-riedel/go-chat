import { Injectable } from '@angular/core';
import { catchError, EMPTY, Subject, switchAll, take, tap } from 'rxjs';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { environment } from '../../../../environments/environment';
import { Event, EventType } from 'src/app/core/model/Chat';

export const WS_ENDPOINT = environment.endpoints.ws;

@Injectable({
  providedIn: 'root'
})
export class ChatService {

  private subject$ : WebSocketSubject<Event> | undefined;
  private messageSubject$ = new Subject<Event>();

  public confConnection(userName : string){
    this.subject$ = webSocket<Event>(WS_ENDPOINT);
    this.subject$?.next({
      type : EventType.SetUsernameEvent,
      data : {
        client : userName,
        message : '',
        timestamp : null
      }
    });
  }

  constructor() {
    //this.subject$ = webSocket<Event>(WS_ENDPOINT);
  }

  public getWS(){
    return this.subject$;
  }

  send(msg : Event){
    this.subject$?.next(msg);
  }

  close(){
    this.subject$?.complete();
  }
}
