import { notifications } from "@mantine/notifications";

export function toastInfo(message, title = "Success") {
  notifications.show({
    title: title,
    message: message,
    color: "success",
    positions: "top-right",
    autoClose: 3000,
  });
}

export function toastWarning(message, title = "Attention !") {
  notifications.show({
    title: title,
    message: message,
    color: "warning",
    positions: "top-right",
    autoClose: 4000,
  });
}

export function toastError(message, title = "Error !!!") {
  notifications.show({
    title: title,
    message: message,
    color: "danger",
    positions: "top-right",
    autoClose: 5000,
  });
}
