// package: tribecart.v1
// file: tribecart/v1/orders.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class Order extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  clearProductIdsList(): void;
  getProductIdsList(): Array<string>;
  setProductIdsList(value: Array<string>): void;
  addProductIds(value: string, index?: number): string;

  getTotalPrice(): number;
  setTotalPrice(value: number): void;

  hasCreatedAt(): boolean;
  clearCreatedAt(): void;
  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Order.AsObject;
  static toObject(includeInstance: boolean, msg: Order): Order.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Order, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Order;
  static deserializeBinaryFromReader(message: Order, reader: jspb.BinaryReader): Order;
}

export namespace Order {
  export type AsObject = {
    id: string,
    userId: string,
    productIdsList: Array<string>,
    totalPrice: number,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class CreateOrderRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  clearProductIdsList(): void;
  getProductIdsList(): Array<string>;
  setProductIdsList(value: Array<string>): void;
  addProductIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrderRequest): CreateOrderRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateOrderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrderRequest;
  static deserializeBinaryFromReader(message: CreateOrderRequest, reader: jspb.BinaryReader): CreateOrderRequest;
}

export namespace CreateOrderRequest {
  export type AsObject = {
    userId: string,
    productIdsList: Array<string>,
  }
}

export class CreateOrderResponse extends jspb.Message {
  hasOrder(): boolean;
  clearOrder(): void;
  getOrder(): Order | undefined;
  setOrder(value?: Order): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrderResponse): CreateOrderResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateOrderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrderResponse;
  static deserializeBinaryFromReader(message: CreateOrderResponse, reader: jspb.BinaryReader): CreateOrderResponse;
}

export namespace CreateOrderResponse {
  export type AsObject = {
    order?: Order.AsObject,
  }
}

export class GetOrderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrderRequest): GetOrderRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOrderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrderRequest;
  static deserializeBinaryFromReader(message: GetOrderRequest, reader: jspb.BinaryReader): GetOrderRequest;
}

export namespace GetOrderRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetOrderResponse extends jspb.Message {
  hasOrder(): boolean;
  clearOrder(): void;
  getOrder(): Order | undefined;
  setOrder(value?: Order): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrderResponse): GetOrderResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOrderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrderResponse;
  static deserializeBinaryFromReader(message: GetOrderResponse, reader: jspb.BinaryReader): GetOrderResponse;
}

export namespace GetOrderResponse {
  export type AsObject = {
    order?: Order.AsObject,
  }
}

