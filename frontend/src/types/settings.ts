/* tslint:disable */
/**
 * This file was automatically generated by json-schema-to-typescript.
 * DO NOT MODIFY IT BY HAND. Instead, modify the source JSONSchema file,
 * and run json-schema-to-typescript to regenerate this file.
 */

/**
 * The version of these settings. Allowed only to be set to 2 with this schema.
 */
export type Version = number;
/**
 * The name of the serial port to which the host should connect.
 */
export type SerialPort = string;
/**
 * The baud rate at which the host should connect to the printer.
 */
export type Baudrate = number;
/**
 * How many bytes of data transmitted by the printer should be kept in-memory.
 */
export type ScrollbackBufferSize = number;
export type Parity = "NONE" | "ODD" | "EVEN";
export type DataBits = number;
/**
 * The path to a folder where biedaprint stores its data (gcode files, updates etc.)
 */
export type DataPath = string;
/**
 * A command which will run when biedaprint starts
 */
export type StartupCommand = string;
export type PresetName = string;
export type HotendTemperature = number;
export type HotbedTemperature = number;
export type FanSpeed = number;
/**
 * Handy temperature presets shown in a dropdown when manually setting the temeperatures.
 */
export type TemperaturePresets = Presets[];

export interface TheBiedaprintSettingsSchema {
  __v: Version;
  serial: Serial;
  general: General;
  temperatures: Temperatures;
  [k: string]: any;
}
/**
 * Describes how the host should connect to the printer via serial.
 */
export interface Serial {
  serialPort: SerialPort;
  baudRate: Baudrate;
  scrollbackBufferSize: ScrollbackBufferSize;
  parity: Parity;
  dataBits: DataBits;
  [k: string]: any;
}
export interface General {
  dataPath: DataPath;
  startupCommand: StartupCommand;
  [k: string]: any;
}
export interface Temperatures {
  temperaturePresets: TemperaturePresets;
  [k: string]: any;
}
export interface Presets {
  name: PresetName;
  hotendTemperature: HotendTemperature;
  hotbedTemperature: HotbedTemperature;
  fanSpeed?: FanSpeed;
  [k: string]: any;
}
