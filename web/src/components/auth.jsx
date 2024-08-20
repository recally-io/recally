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
import { useState } from "react";
import { toastInfo } from "../libs/alert";
import { post } from "../libs/api";

const url = new URL(window.location.href);
const redirect = url.searchParams.get("redirect");

export function AuthenticationForm() {
  const [type, toggle] = useToggle(["login", "register"]);
  const [errMessage, setErrMessage] = useState("");
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

  const register = async () => {
    try {
      const res = await post("/api/v1/auth/register", null, form.values);
      const user = res.data;
      toastInfo(
        "You have successfully registered: " + user.email + "!",
        "Registration successful",
      );
      setTimeout(() => {
        window.location.href = redirect || "/";
      }, 1000);
    } catch (error) {
      setErrMessage(error.message);
    }
  };

  const login = async () => {
    try {
      const res = await post("/api/v1/auth/login", null, form.values);
      const user = res.data;
      toastInfo(
        `You have successfully logged in as ${user.email}`,
        "Login successful",
      );

      setTimeout(() => {
        window.location.href = redirect || "/";
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
