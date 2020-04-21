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
		return <option value={i}>{item.Title}</option>
	}
	
	showCategories(e) {
		if (e.target.value == -1) return;
		this.setState({
			showCategories: true,
			rubricIndex: e.target.value,
			disabled: false	//Включаем кнопку "Удалить рубрику"
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
			return <Faq rubricIndex={this.state.rubricIndex} categoryIndex={this.state.categoryIndex} showQuestions={this.state.showQuestions} />
		}
		return null
	}
	
	addRubric() {
		let newRubric = prompt('Укажите наименование рубрики');
		fetch('/addRubric', {
			method: 'POST',
			body: new URLSearchParams({
				rubricName: newRubric
			})
		})
		.then(resp => resp.text())
		.then(result => {
			if (/^[0-9A-F]{24}$/i.test(result)) {
				data.unshift({
					"Id": result,
					"Title": newRubric,
					"Categories": []
				});
				this.forceUpdate();
			} else {
				alert(result);
			}
		}, err => alert(err))
	}
	
	deleteRubric() {
		fetch('/deleteRubric', {
			method: 'POST',
			body: new URLSearchParams({
				rubricId: data[this.state.rubricIndex].Id
			})
		})
		.then(resp => resp.text())
		.then(err => {
			if (err == '<nil>') {
				alert("Рубрика успешно удалена!");
				data.splice(this.state.rubricIndex, 1);
				this.setState({showCategories: false});
			} else {
				alert(err);
			}
		}, err => alert(err))
	}
	
	editRubric() {
		let newName = prompt('Укажите наименование категории или оставьте поле пустым', data[this.state.rubricIndex].Title);
		fetch('/editRubric', {
			method: 'POST',
			body: new URLSearchParams({
				rubricId: data[this.state.rubricIndex].Id,
				rubricName: newName
			})
		})
		.then(resp => resp.text())
		.then(err => {
			if (err == '<nil>') {
				data[this.state.rubricIndex].Title = newName;
				this.forceUpdate();
			} else {
				alert(err);
			}
		}, err => alert(err))
	}
	

	render() {
		return (
			<div className="container">
				<div>
					<div>Выберите рубрику</div>
					<select onClick={this.showCategories}>
						<option value={-1} hidden selected>-- Не выбрано --</option>
						{data.map(this.createSelect)}
					</select>
					<button className="add" onClick={this.addRubric}>Добавить новую рубрику</button>
					<button className="add" onClick={this.editRubric} disabled={this.state.disabled}>Переименовать рубрику</button>
					<button className="add" onClick={this.deleteRubric} disabled={this.state.disabled}>Удалить рубрику</button>
				</div>
				<Categories rubricIndex={this.state.rubricIndex} showCategories={this.state.showCategories} showQuestions={this.showQuestions} />
				{this.createFaqBlock()}
			</div>
		)
	}
}

class Categories extends React.Component {
	constructor() {
		super();
		this.state = {
			disabled: true,
		}
		this.addCategory = this.addCategory.bind(this);
		this.deleteCategory = this.deleteCategory.bind(this);
		this.editCategory = this.editCategory.bind(this);
		this.chooseCategory = this.chooseCategory.bind(this);
	}
	createSelect(item, i) {
		let title = item.Title;
		if (title == "") title = "Без названия";
		return <option value={i}>{title}</option>
	}
	
	addCategory() {
		let newCategory = prompt('Укажите наименование категории или оставьте поле пустым');
		fetch('/addCategory', {
			method: 'POST',
			body: new URLSearchParams({
				rubricId: data[this.props.rubricIndex].Id,
				categoryName: newCategory
			})
		})
		.then(resp => resp.text())
		.then(result => {
			if (/^[0-9A-F]{24}$/i.test(result)) {
				data[this.props.rubricIndex].Categories.unshift({
					"Id": result,
					"Title": newCategory,
					"Qas": []
				});
				this.forceUpdate();
			} else {
				alert(result);
			}
		}, err => alert(err))
	}
	
	deleteCategory() {
		fetch('/deleteCategory', {
			method: 'POST',
			body: new URLSearchParams({
				categoryId: this.state.currentCategoryId
			})
		})
		.then(resp => resp.text())
		.then(err => {
			if (err == '<nil>') {
				alert("Категория успешно удалена!");
				data[this.props.rubricIndex].Categories.splice(this.state.currentCategoryIndex, 1);
				this.forceUpdate();
			} else {
				alert(err);
			}
		}, err => alert(err))
	}
	
	editCategory() {
		let newName = prompt('Укажите наименование категории или оставьте поле пустым', data[this.props.rubricIndex].Categories[this.state.currentCategoryIndex].Title);
		fetch('/editCategory', {
			method: 'POST',
			body: new URLSearchParams({
				categoryId: this.state.currentCategoryId,
				categoryName: newName
			})
		})
		.then(resp => resp.text())
		.then(err => {
			if (err == '<nil>') {
				data[this.props.rubricIndex].Categories[this.state.currentCategoryIndex].Title = newName;
				this.forceUpdate();
			} else {
				alert(err);
			}
		}, err => alert(err))
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
			return (
				<div>
					<div>Выберите категорию</div>
					<select onClick={this.chooseCategory}>
						<option value={-1} hidden selected>-- Не выбрано --</option>
						{data[this.props.rubricIndex].Categories.map(this.createSelect)}
					</select>
					<button className="add" onClick={this.addCategory}>Добавить новую категорию</button>
					<button className="add" onClick={this.editCategory} disabled={this.state.disabled}>Переименовать категорию</button>
					<button className="add" onClick={this.deleteCategory} disabled={this.state.disabled}>Удалить категорию</button>
				</div>
			)
		}
		return null
	}
}

class Faq extends React.Component {
	constructor(props) {
		super();
		this.state = {
			qas: data[props.rubricIndex].Categories[props.categoryIndex].Qas
		}
		this.createBlocks = this.createBlocks.bind(this);
		this.addQA = this.addQA.bind(this);
		this.deleteQA = this.deleteQA.bind(this);
	}
	
	createBlocks(item, i) {
		return (
			<QABlock itemId={item.Id} index={i} question={item.Question} answer={item.Answer} delete={this.deleteQA}/>
		)
	}
	
	addQA() {
		fetch('/addQA', {
			method: 'POST',
			body: new URLSearchParams({
				categoryId: data[this.props.rubricIndex].Categories[this.props.categoryIndex].Id
			})
		})
		.then(resp => resp.text())
		.then(result => {
			//Если сервер вернул новый идентификатор
			if (/^[0-9A-F]{24}$/i.test(result)) {
				data[this.props.rubricIndex].Categories[this.props.categoryIndex].Qas.unshift({
					"Id": result,
					"Question": "",
					"Answer": ""
				});
				this.forceUpdate();
				data[this.props.rubricIndex].Categories[this.props.categoryIndex].Qas[0].Question = "\n";
				data[this.props.rubricIndex].Categories[this.props.categoryIndex].Qas[0].Answer = "\n";
				this.forceUpdate();
			} else {
				alert(result);
			}
		}, err => alert(err))
	}
	
	deleteQA(qaId, qaIndex) {
		fetch('/deleteQA', {
			method: 'POST',
			body: new URLSearchParams({
				qaId: qaId
			})
		})
		.then(resp => resp.text())
		.then(err => {
			if (err == '<nil>') {
				data[this.props.rubricIndex].Categories[this.props.categoryIndex].Qas.splice(qaIndex, 1)
				this.forceUpdate()
			} else {
				alert(err);
			}
		}, err => alert(err))
	}

	render() {
		return (
			<div>
				<button onClick={this.addQA}>Добавить новый вопрос/ответ</button>
				{data[this.props.rubricIndex].Categories[this.props.categoryIndex].Qas.map(this.createBlocks)}
			</div>
		)
	}
}

class QABlock extends React.Component {
	constructor() {
		super();
		this.state = {
			disabled: true,
			display: "hidden"
		}
		this.update = this.update.bind(this);
	}
	
	update(e) {
		const el = e.target.dataset.elem;
		this.setState({
			[e.target.dataset.name]: e.target.innerText,
			[e.target.dataset.sibling]: e.target[el].innerText,
			disabled: false,
			display: "hidden"
		})
		let question = e.target.innerText,
			answer = e.target[el].innerText;
		if (e.target.dataset.name == "answer") {
			question = e.target[el].innerText;
			answer = e.target.innerText;
		}
		data[0].Categories[0].Qas[this.props.index] = {
			Id: this.props.itemId,
			Question: question,
			Answer: answer
		}
		//this.forceUpdate();
	}
	
	save(id, index) {
		fetch('/save', {
			method: 'POST',
			body: new URLSearchParams({
				id: id,
				index: index,
				question: this.state.question,
				answer: this.state.answer
			})
		})
		.then(resp => resp.text())
		.then(err => {
			if (err == '<nil>') {
				this.setState({
					disabled: true,		//Отключаем кнопку
					display: "visible"	//Показываем уведомление "Сохранено"
				});
			} else {
				alert(err);
			}
		}, err => alert(err))
	}
	
	render() {
		return (
			<div className="faq">
				<div data-name="question" data-sibling="answer" data-elem="nextElementSibling" contentEditable={true} onInput={this.update}>{this.props.question}</div>
				<div data-name="answer" data-sibling="question" data-elem="previousElementSibling" contentEditable={true} onInput={this.update}>{this.props.answer}</div>
				<span>ID: {this.props.itemId}, index: {this.props.index} </span>
				<button onClick={() => this.save(this.props.itemId, this.props.index)} disabled={this.state.disabled}>Сохранить</button>
				<button onClick={() => this.props.delete(this.props.itemId, this.props.index)}>Удалить</button>
				<span className={"notification "+this.state.display}> Сохранено</span>
			</div>
		)
	}
}

const app = document.getElementById('root');
ReactDOM.render(
	<Body />,
	app
)
