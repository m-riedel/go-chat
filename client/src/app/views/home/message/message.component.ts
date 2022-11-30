import { Component, Input, OnInit } from '@angular/core';
import { Event } from 'src/app/core/model/Chat';

@Component({
  selector: 'app-message',
  templateUrl: './message.component.html',
  styleUrls: ['./message.component.scss']
})
export class MessageComponent implements OnInit {

  @Input() event : Event | undefined;

  @Input() username : string = '';

  constructor() { }

  ngOnInit(): void {
  }

}
