package docker

import "testing"

func TestParseDiskUsage(t *testing.T) {
	// 表驱动测试：一组用例，一个结构体代表一个 case
	tests := []struct {
		name      string // 用例名字
		input     string // 输入
		wantType  string // 期望解析出的 Type
		wantSize  string // 期望解析出的 Size
		wantError bool   // 期望是否报错
	}{
		{
			name:     "正常镜像行",
			input:    `{"Type":"Images","TotalCount":"15","Active":"12","Size":"5.012GB","Reclaimable":"1.78GB (35%)"}`,
			wantType: "Images",
			wantSize: "5.012GB",
		},
		{
			name:     "正常 Build Cache 行",
			input:    `{"Type":"Build Cache","TotalCount":"76","Active":"0","Size":"2.955GB","Reclaimable":"2.818GB"}`,
			wantType: "Build Cache",
			wantSize: "2.955GB",
		},
		{
			name:      "非法 JSON 应该报错",
			input:     `{这不是JSON`,
			wantError: true,
		},
	}

	// 遍历每个用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDiskUsage(tt.input)

			// 该报错的场景
			if tt.wantError {
				if err == nil {
					t.Errorf("期望报错，但没报错")
				}
				return
			}

			// 不该报错的场景
			if err != nil {
				t.Fatalf("意外报错: %v", err)
			}
			if got.Type != tt.wantType {
				t.Errorf("Type = %q, 期望 %q", got.Type, tt.wantType)
			}
			if got.Size != tt.wantSize {
				t.Errorf("Size = %q, 期望 %q", got.Size, tt.wantSize)
			}
		})
	}
}
