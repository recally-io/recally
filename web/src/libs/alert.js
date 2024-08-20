import { notifications } from "@mantine/notifications";

export function toastInfo(message, title = "Success") {
  notifications.show({
    title: title,
    message: message,
    color: "green",
    positions: "top-right",
    autoClose: 1000,
  });
}

export function toastWarning(message, title = "Attention !") {
  notifications.show({
    title: title,
    message: message,
    color: "yellow",
    positions: "top-right",
    autoClose: 1000,
  });
}

export function toastError(message, title = "Error !!!") {
  notifications.show({
    title: title,
    message: message,
    color: "red",
    positions: "top-right",
    autoClose: 3000,
  });
}
