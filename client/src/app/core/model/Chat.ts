export interface Event{
    type : EventType;
    data : MessageEventData;
}
export interface MessageEventData{
    client : string;
    message : string;
    timestamp : string | null;
}

export enum EventType{
    JoinRoomEvent, LeaveRoomEvent, MessageEvent, SetUsernameEvent
}