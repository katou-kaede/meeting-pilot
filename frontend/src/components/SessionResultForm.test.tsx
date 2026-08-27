import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  describe,
  expect,
  test,
  vi,
} from "vitest";

import SessionResultForm from "./SessionResultForm";

describe("SessionResultForm", () => {
  test("編集権限がある場合は入力欄を表示する", () => {
    render(
      <SessionResultForm
        decisions="採用する"
        todo="仕様を確認する"
        canEditSession={true}
        onDecisionsChange={vi.fn()}
        onTodoChange={vi.fn()}
      />
    );

    expect(
      screen.getByRole("textbox", {
        name: "決定事項",
      })
    ).toBeInTheDocument();

    expect(
      screen.getByRole("textbox", {
        name: "TODO",
      })
    ).toBeInTheDocument();
  });

  test("編集権限がない場合は閲覧表示にする", () => {
    render(
      <SessionResultForm
        decisions="採用する"
        todo="仕様を確認する"
        canEditSession={false}
        onDecisionsChange={vi.fn()}
        onTodoChange={vi.fn()}
      />
    );

    expect(
      screen.queryByRole("textbox")
    ).not.toBeInTheDocument();

    expect(
      screen.getByText("採用する")
    ).toBeInTheDocument();

    expect(
      screen.getByText("仕様を確認する")
    ).toBeInTheDocument();
  });

  test("決定事項を入力すると変更処理を呼ぶ", async () => {
    const user = userEvent.setup();
    const onDecisionsChange = vi.fn();

    render(
      <SessionResultForm
        decisions=""
        todo=""
        canEditSession={true}
        onDecisionsChange={onDecisionsChange}
        onTodoChange={vi.fn()}
      />
    );

    await user.type(
      screen.getByRole("textbox", {
        name: "決定事項",
      }),
      "採用"
    );

    expect(
      onDecisionsChange
    ).toHaveBeenCalled();
  });
});