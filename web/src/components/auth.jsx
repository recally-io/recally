import { Icon } from "@iconify/react/dist/iconify.js";
import {
  Anchor,
  Button,
  Checkbox,
  Container,
  Divider,
  Group,
  Notification,
  Paper,
  PasswordInput,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { upperFirst, useToggle } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import useStore from "../libs/store";
import { AuthApi } from "../sdk/index";

export function AuthenticationForm() {
  const [type, toggle] = useToggle(["login", "register"]);
  const [errMessage, setErrMessage] = useState("");
  const setUser = useStore((state) => state.setUser);
  const form = useForm({
    initialValues: {
      email: "",
      name: "",
      password: "",
      terms: true,
    },

    validate: {
      email: (val) => (/^\S+@\S+$/.test(val) ? null : "Invalid email"),
      password: (val) =>
        val.length <= 6
          ? "Password should include at least 6 characters"
          : null,
    },
  });

  const authApi = new AuthApi();

  const register = async () => {
    try {
      const user = await authApi.authRegisterPost({
        request: form.values,
      });
      setUser(user);
      console.log(user);
      notifications.show({
        title: "Registration successful",
        message: "You have successfully registered: " + user.email + "!",
        color: "green",
        positions: "top-right",
        autoClose: 1000,
      });
      // redirect to home page
      window.location.href = "/";
    } catch (error) {
      setErrMessage(error.message);
    }
  };

  const login = async () => {
    try {
      const user = await authApi.authLoginPost({ request: form.values });
      setUser(user);
      notifications.show({
        title: "Login successful",
        message: "You have successfully logged in!",
        color: "green",
        positions: "top-right",
        autoClose: 1000,
      });
      // redirect to home page after successful login and wait for 1 second
      setTimeout(() => {
        console.log("Redirecting to home page");
        window.location.href = "/";
      }, 1000);
    } catch (error) {
      setErrMessage(error.message);
    }
  };

  return (
    <>
      <Container size="sm" my="lg">
        {errMessage && (
          <Notification
            color="red"
            title="Login Error!"
            onClose={() => {
              setErrMessage("");
            }}
          >
            {errMessage}
          </Notification>
        )}

        <Paper radius="md" p="xl" withBorder>
          <Text size="lg" fw={500}>
            Welcome to Vibrainb, {type} with
          </Text>

          <Group grow mb="md" mt="md">
            <Button leftSection={<Icon tabler:brand-google />} variant="filled">
              Google
            </Button>
            <Button
              leftSection={<Icon tabler:brand-twitter-filled />}
              variant="filled"
            >
              Twitter
            </Button>
          </Group>

          <Divider
            label="Or continue with email"
            labelPosition="center"
            my="lg"
          />

          <form
            onSubmit={form.onSubmit(() => {
              if (type === "register") {
                register();
              } else {
                login();
              }
            })}
          >
            <Stack>
              {type === "register" && (
                <TextInput
                  label="Name"
                  placeholder="Your name"
                  value={form.values.name}
                  onChange={(event) =>
                    form.setFieldValue("name", event.currentTarget.value)
                  }
                  radius="md"
                />
              )}

              <TextInput
                required
                label="Email"
                placeholder="hello@mantine.dev"
                value={form.values.email}
                onChange={(event) =>
                  form.setFieldValue("email", event.currentTarget.value)
                }
                error={form.errors.email && "Invalid email"}
                radius="md"
              />

              <PasswordInput
                required
                label="Password"
                placeholder="Your password"
                value={form.values.password}
                onChange={(event) =>
                  form.setFieldValue("password", event.currentTarget.value)
                }
                error={
                  form.errors.password &&
                  "Password should include at least 6 characters"
                }
                radius="md"
              />

              {type === "register" && (
                <Checkbox
                  label="I accept terms and conditions"
                  checked={form.values.terms}
                  onChange={(event) =>
                    form.setFieldValue("terms", event.currentTarget.checked)
                  }
                />
              )}
            </Stack>

            <Group justify="space-between" mt="xl">
              <Anchor
                component="button"
                type="button"
                c="dimmed"
                onClick={() => toggle()}
                size="xs"
              >
                {type === "register"
                  ? "Already have an account? Login"
                  : "Don't have an account? Register"}
              </Anchor>
              <Button type="submit" radius="xl">
                {upperFirst(type)}
              </Button>
            </Group>
          </form>
        </Paper>
      </Container>
    </>
  );
}
