// package: tribecart.v1
// file: tribecart/v1/payments.proto

import * as jspb from "google-protobuf";

export class ProcessPaymentRequest extends jspb.Message {
  getOrderId(): string;
  setOrderId(value: string): void;

  getAmount(): number;
  setAmount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProcessPaymentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProcessPaymentRequest): ProcessPaymentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProcessPaymentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProcessPaymentRequest;
  static deserializeBinaryFromReader(message: ProcessPaymentRequest, reader: jspb.BinaryReader): ProcessPaymentRequest;
}

export namespace ProcessPaymentRequest {
  export type AsObject = {
    orderId: string,
    amount: number,
  }
}

export class ProcessPaymentResponse extends jspb.Message {
  getTransactionId(): string;
  setTransactionId(value: string): void;

  getSuccess(): boolean;
  setSuccess(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProcessPaymentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ProcessPaymentResponse): ProcessPaymentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProcessPaymentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProcessPaymentResponse;
  static deserializeBinaryFromReader(message: ProcessPaymentResponse, reader: jspb.BinaryReader): ProcessPaymentResponse;
}

export namespace ProcessPaymentResponse {
  export type AsObject = {
    transactionId: string,
    success: boolean,
  }
}

