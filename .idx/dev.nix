{ pkgs, ... }: {
  channel = "stable-23.11";

  packages = [
    pkgs.go
  ];

  env = { };

  idx = {
    extensions = [
      "golang.go"
    ];
  };
}
