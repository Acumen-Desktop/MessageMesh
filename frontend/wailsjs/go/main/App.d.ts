// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {models} from '../models';

export function GetAccounts():Promise<Array<models.Account>>;

export function GetBlockchain():Promise<Array<models.Block>>;

export function GetMessages():Promise<Array<models.Message>>;

export function GetMessagesFromPeer(arg1:string):Promise<Array<models.Message>>;

export function SendMessage(arg1:string,arg2:string):Promise<void>;
