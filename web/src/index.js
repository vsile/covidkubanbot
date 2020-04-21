class Body extends React.Component {
	constructor() {
		super();
		this.state = {
			showCategories: false,
			showQuestions: false,
			rubricIndex: -1,
			categoryIndex: -1,
			disabled: true
		};
		this.showCategories = this.showCategories.bind(this);
		this.showQuestions = this.showQuestions.bind(this);
		this.addRubric = this.addRubric.bind(this);
		this.deleteRubric = this.deleteRubric.bind(this);
		this.editRubric = this.editRubric.bind(this);
	}
	createSelect(item, i) {
		return React.createElement(
			'option',
			{ value: i },
			item.Title
		);
	}

	showCategories(e) {
		if (e.target.value == -1) return;
		this.setState({
			showCategories: true,
			rubricIndex: e.target.value,
			disabled: false //Включаем кнопку "Удалить рубрику"
		});
	}
	showQuestions(e) {
		this.setState({
			showQuestions: true,
			categoryIndex: e.target.value
		});
	}

	createFaqBlock() {
		if (this.state.showQuestions && this.state.categoryIndex < data[this.state.rubricIndex].Categories.length) {
			return React.createElement(Faq, { rubricIndex: this.state.rubricIndex, categoryIndex: this.state.categoryIndex, showQuestions: this.state.showQuestions });
		}
		return null;
	}

	addRubric() {
		var _this = this;

		var newRubric = prompt('Укажите наименование рубрики');
		fetch('/addRubric', {
			method: 'POST',
			body: new URLSearchParams({
				rubricName: newRubric
			})
		}).then(function (resp) {
			return resp.text();
		}).then(function (result) {
			if (/^[0-9A-F]{24}$/i.test(result)) {
				data.unshift({
					"Id": result,
					"Title": newRubric,
					"Categories": []
				});
				_this.forceUpdate();
			} else {
				alert(result);
			}
		}, function (err) {
			return alert(err);
		});
	}

	deleteRubric() {
		var _this2 = this;

		fetch('/deleteRubric', {
			method: 'POST',
			body: new URLSearchParams({
				rubricId: data[this.state.rubricIndex].Id
			})
		}).then(function (resp) {
			return resp.text();
		}).then(function (err) {
			if (err == '<nil>') {
				alert("Рубрика успешно удалена!");
				data.splice(_this2.state.rubricIndex, 1);
				_this2.setState({ showCategories: false });
			} else {
				alert(err);
			}
		}, function (err) {
			return alert(err);
		});
	}

	editRubric() {
		var _this3 = this;

		var newName = prompt('Укажите наименование категории или оставьте поле пустым', data[this.state.rubricIndex].Title);
		fetch('/editRubric', {
			method: 'POST',
			body: new URLSearchParams({
				rubricId: data[this.state.rubricIndex].Id,
				rubricName: newName
			})
		}).then(function (resp) {
			return resp.text();
		}).then(function (err) {
			if (err == '<nil>') {
				data[_this3.state.rubricIndex].Title = newName;
				_this3.forceUpdate();
			} else {
				alert(err);
			}
		}, function (err) {
			return alert(err);
		});
	}

	render() {
		return React.createElement(
			'div',
			{ className: 'container' },
			React.createElement(
				'div',
				null,
				React.createElement(
					'div',
					null,
					'\u0412\u044B\u0431\u0435\u0440\u0438\u0442\u0435 \u0440\u0443\u0431\u0440\u0438\u043A\u0443'
				),
				React.createElement(
					'select',
					{ onClick: this.showCategories },
					React.createElement(
						'option',
						{ value: -1, hidden: true, selected: true },
						'-- \u041D\u0435 \u0432\u044B\u0431\u0440\u0430\u043D\u043E --'
					),
					data.map(this.createSelect)
				),
				React.createElement(
					'button',
					{ className: 'add', onClick: this.addRubric },
					'\u0414\u043E\u0431\u0430\u0432\u0438\u0442\u044C \u043D\u043E\u0432\u0443\u044E \u0440\u0443\u0431\u0440\u0438\u043A\u0443'
				),
				React.createElement(
					'button',
					{ className: 'add', onClick: this.editRubric, disabled: this.state.disabled },
					'\u041F\u0435\u0440\u0435\u0438\u043C\u0435\u043D\u043E\u0432\u0430\u0442\u044C \u0440\u0443\u0431\u0440\u0438\u043A\u0443'
				),
				React.createElement(
					'button',
					{ className: 'add', onClick: this.deleteRubric, disabled: this.state.disabled },
					'\u0423\u0434\u0430\u043B\u0438\u0442\u044C \u0440\u0443\u0431\u0440\u0438\u043A\u0443'
				)
			),
			React.createElement(Categories, { rubricIndex: this.state.rubricIndex, showCategories: this.state.showCategories, showQuestions: this.showQuestions }),
			this.createFaqBlock()
		);
	}
}

class Categories extends React.Component {
	constructor() {
		super();
		this.state = {
			disabled: true
		};
		this.addCategory = this.addCategory.bind(this);
		this.deleteCategory = this.deleteCategory.bind(this);
		this.editCategory = this.editCategory.bind(this);
		this.chooseCategory = this.chooseCategory.bind(this);
	}
	createSelect(item, i) {
		var title = item.Title;
		if (title == "") title = "Без названия";
		return React.createElement(
			'option',
			{ value: i },
			title
		);
	}

	addCategory() {
		var _this4 = this;

		var newCategory = prompt('Укажите наименование категории или оставьте поле пустым');
		fetch('/addCategory', {
			method: 'POST',
			body: new URLSearchParams({
				rubricId: data[this.props.rubricIndex].Id,
				categoryName: newCategory
			})
		}).then(function (resp) {
			return resp.text();
		}).then(function (result) {
			if (/^[0-9A-F]{24}$/i.test(result)) {
				data[_this4.props.rubricIndex].Categories.unshift({
					"Id": result,
					"Title": newCategory,
					"Qas": []
				});
				_this4.forceUpdate();
			} else {
				alert(result);
			}
		}, function (err) {
			return alert(err);
		});
	}

	deleteCategory() {
		var _this5 = this;

		fetch('/deleteCategory', {
			method: 'POST',
			body: new URLSearchParams({
				categoryId: this.state.currentCategoryId
			})
		}).then(function (resp) {
			return resp.text();
		}).then(function (err) {
			if (err == '<nil>') {
				alert("Категория успешно удалена!");
				data[_this5.props.rubricIndex].Categories.splice(_this5.state.currentCategoryIndex, 1);
				_this5.forceUpdate();
			} else {
				alert(err);
			}
		}, function (err) {
			return alert(err);
		});
	}

	editCategory() {
		var _this6 = this;

		var newName = prompt('Укажите наименование категории или оставьте поле пустым', data[this.props.rubricIndex].Categories[this.state.currentCategoryIndex].Title);
		fetch('/editCategory', {
			method: 'POST',
			body: new URLSearchParams({
				categoryId: this.state.currentCategoryId,
				categoryName: newName
			})
		}).then(function (resp) {
			return resp.text();
		}).then(function (err) {
			if (err == '<nil>') {
				data[_this6.props.rubricIndex].Categories[_this6.state.currentCategoryIndex].Title = newName;
				_this6.forceUpdate();
			} else {
				alert(err);
			}
		}, function (err) {
			return alert(err);
		});
	}

	chooseCategory(e) {
		if (e.target.value == -1) return;
		this.setState({
			disabled: false,
			currentCategoryId: data[this.props.rubricIndex].Categories[e.target.value].Id,
			currentCategoryIndex: e.target.value
		});
		this.props.showQuestions(e);
	}

	render() {
		if (this.props.showCategories) {
			return React.createElement(
				'div',
				null,
				React.createElement(
					'div',
					null,
					'\u0412\u044B\u0431\u0435\u0440\u0438\u0442\u0435 \u043A\u0430\u0442\u0435\u0433\u043E\u0440\u0438\u044E'
				),
				React.createElement(
					'select',
					{ onClick: this.chooseCategory },
					React.createElement(
						'option',
						{ value: -1, hidden: true, selected: true },
						'-- \u041D\u0435 \u0432\u044B\u0431\u0440\u0430\u043D\u043E --'
					),
					data[this.props.rubricIndex].Categories.map(this.createSelect)
				),
				React.createElement(
					'button',
					{ className: 'add', onClick: this.addCategory },
					'\u0414\u043E\u0431\u0430\u0432\u0438\u0442\u044C \u043D\u043E\u0432\u0443\u044E \u043A\u0430\u0442\u0435\u0433\u043E\u0440\u0438\u044E'
				),
				React.createElement(
					'button',
					{ className: 'add', onClick: this.editCategory, disabled: this.state.disabled },
					'\u041F\u0435\u0440\u0435\u0438\u043C\u0435\u043D\u043E\u0432\u0430\u0442\u044C \u043A\u0430\u0442\u0435\u0433\u043E\u0440\u0438\u044E'
				),
				React.createElement(
					'button',
					{ className: 'add', onClick: this.deleteCategory, disabled: this.state.disabled },
					'\u0423\u0434\u0430\u043B\u0438\u0442\u044C \u043A\u0430\u0442\u0435\u0433\u043E\u0440\u0438\u044E'
				)
			);
		}
		return null;
	}
}

class Faq extends React.Component {
	constructor(props) {
		super();
		this.state = {
			qas: data[props.rubricIndex].Categories[props.categoryIndex].Qas
		};
		this.createBlocks = this.createBlocks.bind(this);
		this.addQA = this.addQA.bind(this);
		this.deleteQA = this.deleteQA.bind(this);
	}

	createBlocks(item, i) {
		return React.createElement(QABlock, { itemId: item.Id, index: i, question: item.Question, answer: item.Answer, 'delete': this.deleteQA });
	}

	addQA() {
		var _this7 = this;

		fetch('/addQA', {
			method: 'POST',
			body: new URLSearchParams({
				categoryId: data[this.props.rubricIndex].Categories[this.props.categoryIndex].Id
			})
		}).then(function (resp) {
			return resp.text();
		}).then(function (result) {
			//Если сервер вернул новый идентификатор
			if (/^[0-9A-F]{24}$/i.test(result)) {
				data[_this7.props.rubricIndex].Categories[_this7.props.categoryIndex].Qas.unshift({
					"Id": result,
					"Question": "",
					"Answer": ""
				});
				_this7.forceUpdate();
				data[_this7.props.rubricIndex].Categories[_this7.props.categoryIndex].Qas[0].Question = "\n";
				data[_this7.props.rubricIndex].Categories[_this7.props.categoryIndex].Qas[0].Answer = "\n";
				_this7.forceUpdate();
			} else {
				alert(result);
			}
		}, function (err) {
			return alert(err);
		});
	}

	deleteQA(qaId, qaIndex) {
		var _this8 = this;

		fetch('/deleteQA', {
			method: 'POST',
			body: new URLSearchParams({
				qaId: qaId
			})
		}).then(function (resp) {
			return resp.text();
		}).then(function (err) {
			if (err == '<nil>') {
				data[_this8.props.rubricIndex].Categories[_this8.props.categoryIndex].Qas.splice(qaIndex, 1);
				_this8.forceUpdate();
			} else {
				alert(err);
			}
		}, function (err) {
			return alert(err);
		});
	}

	render() {
		return React.createElement(
			'div',
			null,
			React.createElement(
				'button',
				{ onClick: this.addQA },
				'\u0414\u043E\u0431\u0430\u0432\u0438\u0442\u044C \u043D\u043E\u0432\u044B\u0439 \u0432\u043E\u043F\u0440\u043E\u0441/\u043E\u0442\u0432\u0435\u0442'
			),
			data[this.props.rubricIndex].Categories[this.props.categoryIndex].Qas.map(this.createBlocks)
		);
	}
}

class QABlock extends React.Component {
	constructor() {
		super();
		this.state = {
			disabled: true,
			display: "hidden"
		};
		this.update = this.update.bind(this);
	}

	update(e) {
		const el = e.target.dataset.elem;
		this.setState({
			[e.target.dataset.name]: e.target.innerText,
			[e.target.dataset.sibling]: e.target[el].innerText,
			disabled: false,
			display: "hidden"
		});
		var question = e.target.innerText,
		    answer = e.target[el].innerText;
		if (e.target.dataset.name == "answer") {
			question = e.target[el].innerText;
			answer = e.target.innerText;
		}
		data[0].Categories[0].Qas[this.props.index] = {
			Id: this.props.itemId,
			Question: question,
			Answer: answer
			//this.forceUpdate();
		};
	}

	save(id, index) {
		var _this9 = this;

		fetch('/save', {
			method: 'POST',
			body: new URLSearchParams({
				id: id,
				index: index,
				question: this.state.question,
				answer: this.state.answer
			})
		}).then(function (resp) {
			return resp.text();
		}).then(function (err) {
			if (err == '<nil>') {
				_this9.setState({
					disabled: true, //Отключаем кнопку
					display: "visible" //Показываем уведомление "Сохранено"
				});
			} else {
				alert(err);
			}
		}, function (err) {
			return alert(err);
		});
	}

	render() {
		var _this10 = this;

		return React.createElement(
			'div',
			{ className: 'faq' },
			React.createElement(
				'div',
				{ 'data-name': 'question', 'data-sibling': 'answer', 'data-elem': 'nextElementSibling', contentEditable: true, onInput: this.update },
				this.props.question
			),
			React.createElement(
				'div',
				{ 'data-name': 'answer', 'data-sibling': 'question', 'data-elem': 'previousElementSibling', contentEditable: true, onInput: this.update },
				this.props.answer
			),
			React.createElement(
				'span',
				null,
				'ID: ',
				this.props.itemId,
				', index: ',
				this.props.index,
				' '
			),
			React.createElement(
				'button',
				{ onClick: function () {
						return _this10.save(_this10.props.itemId, _this10.props.index);
					}, disabled: this.state.disabled },
				'\u0421\u043E\u0445\u0440\u0430\u043D\u0438\u0442\u044C'
			),
			React.createElement(
				'button',
				{ onClick: function () {
						return _this10.props.delete(_this10.props.itemId, _this10.props.index);
					} },
				'\u0423\u0434\u0430\u043B\u0438\u0442\u044C'
			),
			React.createElement(
				'span',
				{ className: "notification " + this.state.display },
				' \u0421\u043E\u0445\u0440\u0430\u043D\u0435\u043D\u043E'
			)
		);
	}
}

const app = document.getElementById('root');
ReactDOM.render(React.createElement(Body, null), app);
