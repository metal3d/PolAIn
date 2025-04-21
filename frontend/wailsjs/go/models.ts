export namespace main {
	
	export class ModelPresentation {
	    name: string;
	    description: string;
	    provider: string;
	    uncensored?: boolean;
	    reasoning?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ModelPresentation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.provider = source["provider"];
	        this.uncensored = source["uncensored"];
	        this.reasoning = source["reasoning"];
	    }
	}

}

