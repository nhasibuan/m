package crud

type service struct {
	repo1 *repository
}

func NewService(repo *repository) *service {
	return &service{repo1: repo}
}

func (s *service) select1() ([]dual, error) {
	return s.repo1.select1()
}

func (s *service) selectAll() ([]Category, error) {
	return s.repo1.selectAll()
}

func (s *service) insert(data *Category) error {
	return s.repo1.insert(data)
}

func (s *service) selectByID(id int) (*Category, error) {
	return s.repo1.selectByID(id)
}

func (s *service) update(category *Category) error {
	return s.repo1.update(category)
}

func (s *service) delete(id int) error {
	return s.repo1.delete(id)
}
