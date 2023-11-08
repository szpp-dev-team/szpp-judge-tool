package contest

type Config struct {
	/*
	 * NOTE: title, start_at, end_at なども yaml で管理できるようにするか悩んだが
	 * Web 上で変更した時に GitHub でも反映するのがめんどくさいので
	 * コンテストのメタデータは Web 上のみで管理する運用にする
	 */
	Tasks []*Task `yaml:"tasks"` // 問題一覧
}

type Task struct {
	Slug  string `yaml:"slug"`
	Score int    `yaml:"score"`
}
